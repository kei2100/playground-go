package gateway

import (
	"crypto/tls"
	"io"
	"net"
	"sync"
)

// dialer abstracts TLS/PlainText Dial function
type dialer interface {
	Dial(network, addr string) (net.Conn, error)
}

// newDialer selects TLS/PlainText dialer
func newDialer(noTLS bool) dialer {
	if noTLS {
		return new(net.Dialer)
	}
	return new(tlsDialer)
}

// tlsDialer is dialer for TLS connection
type tlsDialer struct{}

// Dial delegates tls.Dial
func (d *tlsDialer) Dial(network, addr string) (net.Conn, error) {
	return tls.Dial(network, addr, &tls.Config{})
}

type dialOptions struct {
	noTLS bool
}

// DialOption is the functional option for dial
type DialOption func(o *dialOptions)

// WithPlainText option
func WithNoTLS() DialOption {
	return func(o *dialOptions) {
		o.noTLS = true
	}
}

// NewRouter returns router for destAddr
func NewRouter(destAddr string, opts ...DialOption) RouteFunc {
	o := new(dialOptions)

	for _, f := range opts {
		f(o)
	}

	d := newDialer(o.noTLS)

	return func(lc net.Conn) {
		rc, err := d.Dial("tcp", destAddr)
		if err != nil {
			panic(err)
		}

		lco, rco := new(sync.Once), new(sync.Once)
		wg := new(sync.WaitGroup)
		wg.Add(2)

		go func() {
			defer lco.Do(func() { lc.Close() })
			defer rco.Do(func() { rc.Close() })
			defer wg.Done()
			io.Copy(rc, lc)
		}()
		go func() {
			defer rco.Do(func() { rc.Close() })
			defer lco.Do(func() { lc.Close() })
			defer wg.Done()
			io.Copy(lc, rc)
		}()

		wg.Wait()
	}
}
