package http

import (
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/andybalholm/brotli"
)

type decompressRoundTripper struct {
	w http.RoundTripper
}

func (d *decompressRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	r := d.w
	if r == nil {
		r = http.DefaultTransport
	}
	res, err := r.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	ce := res.Header.Get("Content-Encoding")
	if len(ce) == 0 {
		return res, err
	}
	// decompress
	// e.g. `Content-Encoding: deflate, gzip` => decompress `gzip` > `deflate`
	var decompressed bool
	encodings := strings.Split(ce, ",")
	for i := len(encodings) - 1; i >= 0; i-- {
		switch encodings[i] {
		case "gzip":
			decompressed = true
			r, err := gzip.NewReader(res.Body)
			if err != nil {
				return nil, fmt.Errorf("create gzip reader: %w", err)
			}
			res.Body = &cascadeReadCloser{readFrom: r, cascade: res.Body}
		case "deflate":
			decompressed = true
			r := flate.NewReader(res.Body)
			res.Body = &cascadeReadCloser{readFrom: r, cascade: res.Body}
		case "br":
			decompressed = true
			r := brotli.NewReader(res.Body)
			res.Body = &cascadeReadCloser{readFrom: io.NopCloser(r), cascade: res.Body}
		case "identity", "":
			// nop
		default:
			return nil, fmt.Errorf("unsuported content encoding %s", encodings[i])
		}
	}
	if !decompressed {
		return res, nil
	}
	res.Header.Del("Content-Encoding")
	res.Header.Del("Content-Length")
	res.ContentLength = -1
	res.Uncompressed = true
	return res, nil
}

type cascadeReadCloser struct {
	readFrom io.ReadCloser
	cascade  io.Closer
}

func (c *cascadeReadCloser) Read(p []byte) (int, error) {
	return c.readFrom.Read(p)
}

func (c *cascadeReadCloser) Close() error {
	rerr := c.readFrom.Close()
	cerr := c.cascade.Close()
	if rerr != nil && cerr != nil {
		return fmt.Errorf("%s: %s", rerr.Error(), cerr.Error())
	}
	if rerr != nil {
		return rerr
	}
	if cerr != nil {
		return cerr
	}
	return nil
}

func TestDecompress(t *testing.T) {
	t.SkipNow()

	cli := http.Client{
		Transport: &decompressRoundTripper{},
	}
	req, _ := http.NewRequest("GET", "https://example.com", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	resp, err := cli.Do(req)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(resp.ContentLength)
	fmt.Println(resp.Header.Get("Content-Encoding"))
	fmt.Println(resp.Header.Get("Content-Length"))
	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(b))
	err = resp.Body.Close()
	if err != nil {
		t.Error(err)
	}
}
