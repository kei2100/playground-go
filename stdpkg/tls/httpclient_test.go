package tls

import (
	"crypto/tls"
	"log"
	"net/http"

	"testing"

	"net/http/httptest"

	"github.com/kei2100/playground-go/util/ioutil"
)

func TestHttpClient(t *testing.T) {
	// ルート認証局を設定する場合
	//
	// roots := x509.NewCertPool()
	// ok := roots.AppendCertsFromPEM([]byte(rootPEM))
	// if !ok {
	//   panic("failed to parse root certificate")
	// }
	//
	// conf := &tls.Config{RootCAs: roots}

	// オレオレ証明書などの警告を無視する場合
	//
	// conf := &tls.Config{InsecureSkipVerify: true}

	sv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))
	defer sv.Close()

	conf := &tls.Config{}
	tran := &http.Transport{TLSClientConfig: conf}
	client := &http.Client{Transport: tran}

	res, err := client.Get(sv.URL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	b := ioutil.MustReadAll(res.Body)
	log.Println(string(b))
}
