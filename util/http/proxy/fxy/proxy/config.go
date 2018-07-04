package proxy

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/packer/builder/azure/pkcs12"
)

// TLSClientConfig is configuration of TLS client authentication
type TLSClientConfig struct {
	CACertPath string

	PKCS12Path     string
	PKCS12Password string

	caCertPEM []byte
	certPEM   []byte
	keyPEM    []byte
}

// Load cert files given configuration
func (c *TLSClientConfig) Load() error {
	loader := new(errorOrLoader)
	loader.load(c.loadCACert)
	loader.load(c.loadPKCS12)
	return loader.Err()
}

type errorOrLoader struct {
	err error
}

func (e *errorOrLoader) load(loader func() error) {
	if e.err != nil {
		return
	}
	e.err = loader()
}

func (e *errorOrLoader) Err() error {
	return e.err
}

func (c *TLSClientConfig) loadCACert() error {
	if c.CACertPath == "" {
		return nil
	}
	b, err := ioutil.ReadFile(c.CACertPath)
	if err != nil {
		return fmt.Errorf("proxy: failed to load ca cert file %v : %v", c.CACertPath, err)
	}
	c.caCertPEM = b
	return nil
}

func (c *TLSClientConfig) loadPKCS12() error {
	if c.PKCS12Path == "" {
		return nil
	}
	b, err := ioutil.ReadFile(c.PKCS12Path)
	if err != nil {
		return fmt.Errorf("proxy: failed to load pkcs12 file %v : %v", c.PKCS12Path, err)
	}
	key, cert, err := pkcs12.Decode(b, c.PKCS12Password)
	if err != nil {
		return fmt.Errorf("proxy: failed to decode pkcs12 data: %v", err)
	}
	kp, err := encodePrivateKeyPEMToMemory(key)
	if err != nil {
		return err
	}
	c.keyPEM = kp
	c.certPEM = encodeCertPEMToMemory(cert)
	return nil
}

func encodePrivateKeyPEMToMemory(key interface{}) ([]byte, error) {
	switch k := key.(type) {
	case *rsa.PrivateKey:
		kb := x509.MarshalPKCS1PrivateKey(k)
		pemb := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: kb})
		return pemb, nil
	case *ecdsa.PrivateKey:
		kb, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			return nil, fmt.Errorf("proxy: failed to marshal ecdsa private key: %v", err)
		}
		pemb := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: kb})
		return pemb, nil
	default:
		return nil, fmt.Errorf("proxy: unknown private key type %T", key)
	}
}

func encodeCertPEMToMemory(cert *x509.Certificate) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type: "CERTIFICATE", Bytes: cert.Raw,
	})
}
