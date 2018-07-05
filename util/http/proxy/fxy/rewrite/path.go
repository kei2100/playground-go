package rewrite

import "net/url"

// PathRewriter is an interface to rewrite the path
type PathRewriter interface {
	// Do rewrites the URL
	Do(*url.URL) (rewrited bool)
}

// NewRewriter creates the PathRewriter
func NewRewriter(old, new string) (PathRewriter, error) {
	return newRegexpPathRewriter(old, new)
}

// regexpPathRewriter is an implementation of the PathRewriter using regexp
type regexpPathRewriter struct {
}

func newRegexpPathRewriter(re, repl string) (*regexpPathRewriter, error) {
	return nil, nil
}

func (r *regexpPathRewriter) Do(url *url.URL) bool {
	return false
}
