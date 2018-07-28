package proxy

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"crypto/tls"

	"github.com/kei2100/playground-go/util/http/proxy/fxy/errors"
	"github.com/kei2100/playground-go/util/http/proxy/fxy/rewrite"
	"golang.org/x/crypto/pkcs12"
)

// Config for the proxy Server
type Config struct {
	URLConfig
	HeaderConfig
	TLSClientConfig
}

// Load resources by given configuration
func (c *Config) Load() error {
	e := errors.NewMultiLine()
	loaders := []func() error{
		c.URLConfig.load,
		c.TLSClientConfig.load,
	}

	for _, l := range loaders {
		if err := l(); err != nil {
			e.Add(err)
		}
	}
	if e.Len() > 0 {
		return e
	}
	return nil
}

// HeaderConfig is a configuration of the Request Header
type HeaderConfig struct {
	Header http.Header
}

// URLConfig is a configuration of the URL
type URLConfig struct {
	Destination        string            // destination proto://host[:port]
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
	return c.scheme
}

// UserInfo returns the *url.UserInfo or nil
func (c *URLConfig) UserInfo() *url.Userinfo {
	return c.userInfo
}

// PathRewriters returns path rewriters
func (c *URLConfig) PathRewriters() []rewrite.PathRewriter {
	return c.pathRewriters
}

// load url infos given configuration
func (c *URLConfig) load() error {
	if c == nil {
		return nil
	}

	u, err := url.Parse(c.Destination)
	if err != nil {
		return fmt.Errorf("config: failed to parse Destination string to URL: %v", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("config: invalid scheme %v", u.Scheme)
	}
	c.host = u.Host
	c.scheme = u.Scheme

	if c.Username != "" {
		c.userInfo = url.UserPassword(c.Username, c.Password)
	}

	for old, new := range c.RewritePathEntries {
		rwr, err := rewrite.NewRewriter(old, new)
		if err != nil {
			return fmt.Errorf("config: failed to interpret the rewrite string %v to %v", old, new)
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

	tlsConfig *tls.Config
}

// TLSConfig returns *tls.Config for the tls client certification
func (c *TLSClientConfig) TLSConfig() *tls.Config {
	return c.tlsConfig
}

// load cert files given configuration
func (c *TLSClientConfig) load() error {
	if c == nil {
		return nil
	}

	loader := new(errorOrLoader)
	loader.load(c.loadCACert)
	loader.load(c.loadPKCS12)
	loader.load(c.loadTLSConfig)
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
		return fmt.Errorf("config: failed to load ca cert file %v : %v", c.CACertPath, err)
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
		return fmt.Errorf("config: failed to load pkcs12 file %v : %v", c.PKCS12Path, err)
	}
	key, cert, err := pkcs12.Decode(b, c.PKCS12Password)
	if err != nil {
		return fmt.Errorf("config: failed to decode pkcs12 data: %v", err)
	}
	kp, err := encodePrivateKeyPEMToMemory(key)
	if err != nil {
		return err
	}
	c.keyPEM = kp
	c.certPEM = encodeCertPEMToMemory(cert)
	return nil
}

func (c *TLSClientConfig) loadTLSConfig() error {
	cfg := tls.Config{}

	if certPEM := c.certPEM; certPEM != nil {
		cert, err := tls.X509KeyPair(certPEM, c.keyPEM)
		if err != nil {
			return fmt.Errorf("config: failed to create x509 keypair: %v", err)
		}
		cfg.Certificates = []tls.Certificate{cert}
	}
	if certPEM := c.caCertPEM; certPEM != nil {
		p := x509.NewCertPool()
		if ok := p.AppendCertsFromPEM(certPEM); !ok {
			return fmt.Errorf("config: failed to append ca cert file %v (%v bytes)", c.CACertPath, len(certPEM))
		}
		cfg.RootCAs = p
	}

	c.tlsConfig = &cfg
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
			return nil, fmt.Errorf("config: failed to marshal ecdsa private key: %v", err)
		}
		pemb := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: kb})
		return pemb, nil
	default:
		return nil, fmt.Errorf("config: unknown private key type %T", key)
	}
}

func encodeCertPEMToMemory(cert *x509.Certificate) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type: "CERTIFICATE", Bytes: cert.Raw,
	})
}
