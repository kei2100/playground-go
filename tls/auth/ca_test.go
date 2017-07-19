package main

import (
	"testing"
	"crypto/x509"
	"crypto/x509/pkix"
	"net"
	"crypto/rand"
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

	der, err := x509.CreateCertificateRequest(rand.Reader, &template, genRSAPrivateKey())
	if err != nil {
		t.Fatal(err)
	}
	csr, err := x509.ParseCertificateRequest(der)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := ca.Sign(csr); err != nil {
		t.Fatal(err)
	}
}
