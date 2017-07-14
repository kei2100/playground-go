package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"log"
	"math/big"
	"time"
	"encoding/pem"
	"net"
)

func main() {
}

func createServerCert(rootCert *x509.Certificate) (serverCertPEM []byte){
	template := certTemplate()
	template.KeyUsage = x509.KeyUsageDigitalSignature
	template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	template.IPAddresses = []net.IP{net.ParseIP("127.0.0.1")}

	priv := genRSAPrivateKey()
	_, serverCertPEM = createCert(template, rootCert, &priv.PublicKey, priv)
	return
}

func createRootCert() (rootCert *x509.Certificate, rootCertPEM []byte){
	template := certTemplate()
	template.IsCA = true
	template.KeyUsage = x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature
	template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}
	template.IPAddresses = []net.IP{net.ParseIP("127.0.0.1")}

	priv := genRSAPrivateKey()
	rootCert, rootCertPEM = createCert(template, template, &priv.PublicKey, priv)
	return
}

func genRSAPrivateKey() *rsa.PrivateKey {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("failed to gen key: %v", err)
	}
	return privKey
}

func createCert(template, parent *x509.Certificate, pub interface{}, parentPriv interface{}) (
	cert *x509.Certificate, certPEM []byte) {

	certDER, err := x509.CreateCertificate(rand.Reader, template, parent, pub, parentPriv)
	if err != nil {
		log.Fatalf("failed to create cert: %v", err)
	}
	// parse the resulting certificate so we can use it again
	cert, err = x509.ParseCertificate(certDER)
	if err != nil {
		log.Fatalf("failed to parse cert: %v", err)
	}
	// PEM encode the certificate (this is a standard TLS encoding)
	b := pem.Block{Type: "CERTIFICATE", Bytes: certDER}
	certPEM = pem.EncodeToMemory(&b)
	return
}

func certTemplate() *x509.Certificate {
	// generate a random serial number (a real cert authority would have some logic behind this)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("failed to gen serial number: %v", err)
	}

	var _ = pkix.Name{Organization: []string{"test, Inc."}}

	return &x509.Certificate{
		SerialNumber: serialNumber,
		//Subject:               pkix.Name{Organization: []string{"test, Inc."}},
		SignatureAlgorithm:    x509.SHA256WithRSA,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour),
		BasicConstraintsValid: true,
	}
}
