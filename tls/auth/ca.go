package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"log"
	"net"
	"time"
	"math/big"
)

// CA acts as a certificate authority
type CA struct {
	priv *rsa.PrivateKey
	Cert *x509.Certificate
}

// NewRootCA returns new root CA
func NewRootCA() *CA {
	ca := &CA{}
	ca.priv = genRSAPrivateKey()

	// TODO 項目精査
	template := x509.Certificate{
		SerialNumber: randSerialNumber(),
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(24 * time.Hour),

		KeyUsage:    x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		IsCA:        true,
		IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &ca.priv.PublicKey, ca.priv)
	if err != nil {
		log.Fatalf("failed to create root certificate: %v", err)
	}

	ca.Cert, err = x509.ParseCertificate(certDER)
	if err != nil {
		log.Fatalf("failed to parse certDER: %v", err)
	}

	return ca
}

// x509.CreateCertificateRequestで作成したreqを渡してもらう。
//
func (ca *CA) Sign(csr *x509.CertificateRequest) (*x509.Certificate, error) {
	// TODO 項目精査
	template := x509.Certificate{
		Signature:          csr.Signature,
		SignatureAlgorithm: csr.SignatureAlgorithm,

		PublicKeyAlgorithm: csr.PublicKeyAlgorithm,
		PublicKey:          csr.PublicKey,

		SerialNumber: big.NewInt(2),
		Issuer:       ca.Cert.Subject,
		Subject:      csr.Subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &ca.priv.PublicKey, ca.priv)
	if err != nil {
		log.Fatalf("failed to create root certificate: %v", err)
	}

	return x509.ParseCertificate(certDER)
}
