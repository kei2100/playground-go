package pem

import (
	"crypto/rsa"
	"crypto/x509"
	gopem "encoding/pem"
)

// EncodeCertificateToMemory encodes x509 certificate to PEM
func EncodeCertificateToMemory(c *x509.Certificate) []byte {
	return gopem.EncodeToMemory(&gopem.Block{
		Type: "CERTIFICATE", Bytes: c.Raw,
	})
}

// EncodeRSAPrivateKeyToMemory encodes rsa private key to PEM
func EncodeRSAPrivateKeyToMemory(k *rsa.PrivateKey) []byte {
	return gopem.EncodeToMemory(&gopem.Block{
		Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k),
	})
}
