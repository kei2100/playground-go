package x509

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"math/big"
	"time"
)

// TemplateParam is a struct for create certificate template
type TemplateParam struct {
	SerialNumber *big.Int
	Subject      pkix.Name
	PublicKey    crypto.PublicKey
	HashFunction crypto.Hash
	NotBefore    time.Time
	NotAfter     time.Time
}

// CreateCATemplate creates certificate template for CA
func CreateCATemplate(param TemplateParam) (*x509.Certificate, error) {
	tmpl, err := createCommonCertTemplate(param)
	if err != nil {
		return nil, err
	}
	tmpl.KeyUsage = x509.KeyUsageCertSign | x509.KeyUsageCRLSign
	tmpl.BasicConstraintsValid = true
	tmpl.IsCA = true
	return tmpl, nil
}

// CreateServerTemplate creates certificate template for server
func CreateServerTemplate(param TemplateParam, authorityKeyID []byte) (*x509.Certificate, error) {
	tmpl, err := createCommonCertTemplate(param)
	if err != nil {
		return nil, err
	}
	tmpl.AuthorityKeyId = authorityKeyID
	tmpl.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	tmpl.KeyUsage = x509.KeyUsageDigitalSignature
	tmpl.BasicConstraintsValid = true
	tmpl.IsCA = false
	return tmpl, nil
}

// CreateClientTemplate creates certificate template for server
func CreateClientTemplate(param TemplateParam, authorityKeyID []byte) (*x509.Certificate, error) {
	tmpl, err := createCommonCertTemplate(param)
	if err != nil {
		return nil, err
	}
	tmpl.AuthorityKeyId = authorityKeyID
	tmpl.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	tmpl.KeyUsage = x509.KeyUsageDigitalSignature
	tmpl.BasicConstraintsValid = true
	tmpl.IsCA = false
	return tmpl, nil
}

func createCommonCertTemplate(param TemplateParam) (*x509.Certificate, error) {
	subjectKeyID, err := deriveSubjectKeyID(param.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("cert: failed to marshal public key: %v", err)
	}
	sigAlgo, pubAlgo := deriveAlgorithm(param.PublicKey, param.HashFunction)

	tmpl := x509.Certificate{
		SerialNumber: param.SerialNumber,

		SignatureAlgorithm: sigAlgo,
		PublicKeyAlgorithm: pubAlgo,

		Subject: param.Subject,

		NotBefore: param.NotBefore,
		NotAfter:  param.NotAfter,

		SubjectKeyId: subjectKeyID[:],
	}
	return &tmpl, nil
}

// Sign certificate.
// The returned slice is the certificate in DER encoding
func Sign(signer *x509.Certificate, signerPrivate crypto.Signer, signee *x509.Certificate, signeePublic interface{}) (*x509.Certificate, error) {
	certDER, err := x509.CreateCertificate(rand.Reader, signee, signer, signeePublic, signerPrivate)
	if err != nil {
		return nil, fmt.Errorf("cert: failed to sign: %v", err)
	}
	return x509.ParseCertificate(certDER)
}

func deriveSubjectKeyID(pub crypto.PublicKey) ([20]byte, error) {
	var keyBytes []byte
	switch p := pub.(type) {
	case *rsa.PublicKey:
		keyBytes = x509.MarshalPKCS1PublicKey(p)
	case *ecdsa.PublicKey:
		keyBytes = elliptic.Marshal(p.Curve, p.X, p.Y)
	default:
		return [20]byte{}, errors.New("cert: only RSA and ECDSA public keys supported")
	}
	return sha1.Sum(keyBytes), nil
}

func deriveAlgorithm(pub crypto.PublicKey, hf crypto.Hash) (x509.SignatureAlgorithm, x509.PublicKeyAlgorithm) {
	sigAlgo, pubAlgo := x509.UnknownSignatureAlgorithm, x509.UnknownPublicKeyAlgorithm

	switch pub.(type) {
	case *rsa.PublicKey:
		pubAlgo = x509.RSA
		switch hf {
		case crypto.MD5:
			sigAlgo = x509.MD5WithRSA
		case crypto.SHA256:
			sigAlgo = x509.SHA256WithRSA
		case crypto.SHA384:
			sigAlgo = x509.SHA384WithRSA
		case crypto.SHA512:
			sigAlgo = x509.SHA512WithRSA
		}
	case *ecdsa.PublicKey:
		pubAlgo = x509.ECDSA
		switch hf {
		case crypto.SHA256:
			sigAlgo = x509.ECDSAWithSHA256
		case crypto.SHA384:
			sigAlgo = x509.ECDSAWithSHA384
		case crypto.SHA512:
			sigAlgo = x509.ECDSAWithSHA512
		}
	}
	return sigAlgo, pubAlgo
}

// CRLParam is a struct for create CRL
type CRLParam struct {
	Revoked    []pkix.RevokedCertificate
	PrivateKey crypto.Signer
	ThisUpdate time.Time
	NextUpdate time.Time
}

// CreateCRL with given parameter
func CreateCRL(caCert *x509.Certificate, param CRLParam) (*pkix.CertificateList, error) {
	der, err := caCert.CreateCRL(rand.Reader, param.PrivateKey, param.Revoked, param.ThisUpdate, param.NextUpdate)
	if err != nil {
		return nil, fmt.Errorf("cert: failed to create CRL: %v", err)
	}
	return x509.ParseCRL(der)
}
