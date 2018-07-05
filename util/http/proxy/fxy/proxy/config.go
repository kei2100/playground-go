package proxy

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/kei2100/playground-go/util/http/proxy/fxy/rewrite"
	"golang.org/x/crypto/pkcs12"
)

// TODO
//type HeaderConfig struct {
//	ForwardHostHeader bool
//}

// URLConfig is a configuration of the URL
type URLConfig struct {
	Server             string            // destination server. proto://host[:port]
	Username           string            // username or blank
	Password           string            // password or blank
	RewritePathEntries map[string]string // map[oldPath]newPath

	host          string // host or host:port
	scheme        string // http or https
	userInfo      *url.Userinfo
	pathRewriters []rewrite.PathRewriter
}

// Host returns the host or host:port
func (c *URLConfig) Host() string {
	return c.host
}

// Scheme returns the scheme
func (c *URLConfig) Scheme() string {
	return c.host
}

// UserInfo returns the *url.UserInfo or nil
func (c *URLConfig) UserInfo() *url.Userinfo {
	return c.userInfo
}

// PathRewriters returns path rewriters
func (c *URLConfig) PathRewriters() []rewrite.PathRewriter {
	return c.pathRewriters
}

// Load url infos given configuration
func (c *URLConfig) Load() error {
	u, err := url.Parse(c.Server)
	if err != nil {
		return fmt.Errorf("proxy: failed to parse Server string to URL: %v", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("proxy: invalid scheme %v", u.Scheme)
	}
	c.host = u.Host
	c.scheme = u.Scheme

	if c.Username != "" {
		c.userInfo = url.UserPassword(c.Username, c.Password)
	}

	for old, new := range c.RewritePathEntries {
		rwr, err := rewrite.NewRewriter(old, new)
		if err != nil {
			return fmt.Errorf("proxy: failed to interpret the rewrite string %v to %v", old, new)
		}
		c.pathRewriters = append(c.pathRewriters, rwr)
	}
	return nil
}

// TLSClientConfig is configuration of TLS client authentication
type TLSClientConfig struct {
	CACertPath string

	PKCS12Path     string
	PKCS12Password string

	caCertPEM []byte
	certPEM   []byte
	keyPEM    []byte
}

// CACertPEM returns ca certificate pem data
func (c *TLSClientConfig) CACertPEM() []byte {
	return c.caCertPEM
}

// CertPEM returns client certificate pem data
func (c *TLSClientConfig) CertPEM() []byte {
	return c.certPEM
}

// KeyPEM returns private key pem data for client certification
func (c *TLSClientConfig) KeyPEM() []byte {
	return c.keyPEM
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
