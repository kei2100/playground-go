package tls

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
		IsCA: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, pub, priv)
	if err != nil {
		t.Fatal(err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		t.Fatal(err)
	}

	certPEM := pem.EncodeCertificateToMemory(cert)
	// 出力されたpem文字列をファイルにして以下で内容確認
	// openssl x509 -in file.pem -text -noout
	fmt.Println(string(certPEM))
}
