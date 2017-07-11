package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/kei2100/playground-go/util/ioutil"
)

func main() {
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

	conf := &tls.Config{}
	tran := &http.Transport{TLSClientConfig: conf}
	client := &http.Client{Transport: tran}

	res, err := client.Get("https://google.co.jp/robots.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	b := ioutil.MustReadAll(res.Body)
	log.Println(string(b))
}
