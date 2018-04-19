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

	"io/ioutil"

	// 現状、標準パッケージだとpkcs12.Encodeが存在しない
	"github.com/hashicorp/packer/builder/azure/pkcs12"

	"encoding/base64"

	"github.com/kei2100/playground-go/util/encoding/pem"
	"path/filepath"
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

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	certDER, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, priv.Public(), priv)
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

	// サーバー証明書
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
		IsCA: false,
	}

	srvPriv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	srvCertDER, err := x509.CreateCertificate(rand.Reader, &srvTmpl, ca, srvPriv.Public(), cakey)
	if err != nil {
		t.Fatal(err)
	}
	srvCert, err := x509.ParseCertificate(srvCertDER)
	if err != nil {
		t.Fatal(err)
	}

	// クライアント証明書
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
		IsCA: false,
	}

	clientPriv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	clientCertDER, err := x509.CreateCertificate(rand.Reader, &clientTmpl, ca, clientPriv.Public(), cakey)
	if err != nil {
		t.Fatal(err)
	}
	clientCert, err := x509.ParseCertificate(clientCertDER)
	if err != nil {
		t.Fatal(err)
	}
	clientPFX, err := pkcs12.Encode(clientCertDER, clientPriv, "password")
	if err != nil {
		t.Fatal(err)
	}

	// CRL
	r := pkix.RevokedCertificate{
		SerialNumber:   clientCert.SerialNumber,
		RevocationTime: time.Now(),
	}

	crlDER, err := ca.CreateCRL(rand.Reader, cakey, []pkix.RevokedCertificate{r}, time.Now(), ca.NotAfter) // time.Now()=ThisUpdate, ca.NotAfter=NextUpdate
	if err != nil {
		t.Fatal(err)
	}
	crl, err := x509.ParseCRL(crlDER)
	if err != nil {
		t.Fatal(err)
	}

	// 出力
	{
		fmt.Println("######### CA CERT ###########")
		n, b := "cacert.pem", pem.EncodeCertificateToMemory(ca)
		fmt.Println(string(b))
		writeFile(n, b)
	}
	{
		fmt.Println("######### SERVCER CERT ###########")
		n, b := "servcert.pem", pem.EncodeCertificateToMemory(srvCert)
		fmt.Println(string(b))
		writeFile(n, b)
	}
	{
		fmt.Println("######### SERVCER KEY ###########")
		n, b := "servkey.pem", pem.EncodeRSAPrivateKeyToMemory(srvPriv)
		fmt.Println(string(b))
		writeFile(n, b)
	}
	{
		fmt.Println("######### CLIENT PKCS#12/pfx file ###########")
		n, b := "client.pfx", clientPFX
		fmt.Println(base64.StdEncoding.EncodeToString(b))
		writeFile(n, b)
	}
	{
		fmt.Println("######### CRL ###########")
		n, b := "crl.pem", pem.EncodeCRLToMemory(crl)
		fmt.Println(string(b))
		writeFile(n, b)
	}
}

func writeFile(name string, b []byte) {
	// FIXME 実際に書き込みたいときは false > true に
	if false {
		p := filepath.Join(".", name)
		ioutil.WriteFile(p, b, 0644)
	}
}
