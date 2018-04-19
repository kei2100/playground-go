package pki

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/kei2100/playground-go/util/encoding/pem"
)

func TestSelfSignAsCA(t *testing.T) {
	cert, _ := testSelfSignAsCA(t)
	certPEM := pem.EncodeCertificateToMemory(cert)
	// 出力されたpem文字列をファイルにして以下で内容確認
	// openssl x509 -in file.pem -text -noout
	fmt.Println(string(certPEM))
}

func testSelfSignAsCA(t *testing.T) (*x509.Certificate, *rsa.PrivateKey) {
	t.Helper()

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	pub := priv.Public()

	var sn int64 = 1
	var subjKeyID = []byte{1, 2, 3, 4}

	tmpl := x509.Certificate{
		SerialNumber: big.NewInt(sn),

		SignatureAlgorithm: x509.SHA256WithRSA,
		PublicKeyAlgorithm: x509.RSA,

		Subject: pkix.Name{
			Country:      []string{"JP"},
			Organization: []string{"MyCompany, inc."},
			CommonName:   "MyCompanyCA",
		},

		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(10, 0, 0),

		SubjectKeyId:          subjKeyID,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, pub, priv)
	if err != nil {
		t.Fatal(err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		t.Fatal(err)
	}

	return cert, priv
}

func TestIssueCerts(t *testing.T) {
	ca, cakey := testSelfSignAsCA(t)
	ckID := ca.SubjectKeyId

	srvPriv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	srvTmpl := x509.Certificate{
		SerialNumber: big.NewInt(ca.SerialNumber.Int64() + 1),

		SignatureAlgorithm: x509.SHA256WithRSA,
		PublicKeyAlgorithm: x509.RSA,

		Subject: pkix.Name{
			Country:      []string{"JP"},
			Organization: []string{"MyCompany, inc."},
			CommonName:   "mycompany.example.com",
		},

		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(0, 0, 825),

		SubjectKeyId:   []byte{ckID[0], ckID[1], ckID[2], ckID[3] + 1},
		AuthorityKeyId: ca.SubjectKeyId,

		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageDataEncipherment,
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	srvCertDER, err := x509.CreateCertificate(rand.Reader, &srvTmpl, ca, srvPriv.Public(), cakey)
	if err != nil {
		t.Fatal(err)
	}

	srvCert, err := x509.ParseCertificate(srvCertDER)
	if err != nil {
		t.Fatal(err)
	}


	clientTmpl := x509.Certificate{
		SerialNumber: big.NewInt(ca.SerialNumber.Int64() + 2),

		SignatureAlgorithm: x509.SHA256WithRSA,
		PublicKeyAlgorithm: x509.RSA,

		Subject: pkix.Name{
			Country:      []string{"JP"},
			Organization: []string{"MyCompany, inc."},
			CommonName:   "MyClient",
		},

		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(0, 0, 825),

		SubjectKeyId:   []byte{ckID[0], ckID[1], ckID[2], ckID[3] + 2},
		AuthorityKeyId: ca.SubjectKeyId,

		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  false,
	}


	fmt.Println("######### CA CERT ###########")
	fmt.Println(string(pem.EncodeCertificateToMemory(ca)))

	fmt.Println("######### SERVCER CERT ###########")
	fmt.Println(string(pem.EncodeCertificateToMemory(srvCert)))

	fmt.Println("######### SERVCER KEY ###########")
	fmt.Println(string(pem.EncodeRSAPrivateKeyToMemory(srvPriv)))
}
