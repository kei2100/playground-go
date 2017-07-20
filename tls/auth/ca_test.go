package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kei2100/playground-go/util/encoding/pem"
	"fmt"
	"net/http/httputil"
)

func TestNewRootCA(t *testing.T) {
	t.Parallel()

	ca := NewRootCA()

	template := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:   "test.example.com",
			Organization: []string{"Σ Acme Co"},
		},
		SignatureAlgorithm: x509.SHA256WithRSA,
		DNSNames:           []string{"test.example.com"},
		EmailAddresses:     []string{"gopher@golang.org"},
		IPAddresses:        []net.IP{net.IPv4(127, 0, 0, 1).To4(), net.ParseIP("2001:4860:0:2001::68")},
	}

	priv := genRSAPrivateKey()

	der, err := x509.CreateCertificateRequest(rand.Reader, &template, priv)
	if err != nil {
		t.Fatal(err)
	}
	csr, err := x509.ParseCertificateRequest(der)
	if err != nil {
		t.Fatal(err)
	}

	cert, err := ca.Sign(csr)
	if err != nil {
		t.Fatal(err)
	}
	certPEM := pem.EncodeCertificateToMemory(cert)
	privPEM := pem.EncodeRSAPrivateKeyToMemory(priv)

	servTLSCert, err := tls.X509KeyPair(certPEM, privPEM)
	if err != nil {
		log.Fatalf("invalid key pair: %v", err)
	}
	// create another test server and use the certificate
	s := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("HI!")) }))
	s.TLS = &tls.Config{
		Certificates: []tls.Certificate{servTLSCert},
	}

	s.StartTLS()

	// create a pool of trusted certs
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(pem.EncodeCertificateToMemory(ca.Certificate))

	// configure a client to use trust those certificates
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: certPool, InsecureSkipVerify: true},
		},
	}
	resp, err := client.Get(s.URL)
	s.Close()
	if err != nil {
		log.Fatalf("could not make GET request: %v", err)
	}
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatalf("could not dump response: %v", err)
	}
	fmt.Printf("%s\n", dump)
}
