package pem

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	gopem "encoding/pem"
	"fmt"
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

// EncodePrivateKeyToMemory encodes private key to PEM
func EncodePrivateKeyToMemory(key interface{}) ([]byte, error) {
	switch k := key.(type) {
	case *rsa.PrivateKey:
		return EncodeRSAPrivateKeyToMemory(k), nil
	case *ecdsa.PrivateKey:
		kb, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			return nil, fmt.Errorf("pem: failed to marshal ecdsa private key: %v", err)
		}
		pb := gopem.EncodeToMemory(&gopem.Block{Type: "EC PRIVATE KEY", Bytes: kb})
		return pb, nil
	default:
		return nil, fmt.Errorf("pem: unknown private key type %T", key)
	}
}

// EncodeCRLToMemory encodes CRL to PEM
func EncodeCRLToMemory(crl *pkix.CertificateList) []byte {
	raw, _ := asn1.Marshal(*crl) // ignore error
	return gopem.EncodeToMemory(&gopem.Block{
		Type: "X509 CRL", Bytes: raw,
	})
}
