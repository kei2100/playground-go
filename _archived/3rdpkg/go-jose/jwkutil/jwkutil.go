package jwkutil

import (
	"crypto/rand"

	"github.com/square/go-jose/v3"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/xerrors"
)

// Generate JWK pair
// https://openid-foundation-japan.github.io/rfc7517.ja.html
func Generate(use, algorithm, kid string) (public, private string, err error) {
	var pubKey, privKey interface{}

	switch use {
	case "sig":
		switch algorithm {
		case "EdDSA":
			pubKey, privKey, err = ed25519.GenerateKey(rand.Reader)
			if err != nil {
				return "", "", xerrors.Errorf("jwkutil: Generate use `%s` algorithm `%s` generateKey: %w", use, algorithm, err)
			}
		default:
			return "", "", xerrors.Errorf("jwkutil: Generate use `%s` algorithm `%s` not supported", use, algorithm)
		}
	case "enc":
		return "", "", xerrors.New("jwkutil: Generate use `enc` not supported")
	default:
		return "", "", xerrors.Errorf("jwkutil: Generate use `%s` not supported", use)
	}

	pubJWK := jose.JSONWebKey{Key: pubKey, KeyID: kid, Algorithm: algorithm, Use: use}
	privJWK := jose.JSONWebKey{Key: privKey, KeyID: kid, Algorithm: algorithm, Use: use}

	if privJWK.IsPublic() || !pubJWK.IsPublic() || !privJWK.Valid() || !pubJWK.Valid() {
		return "", "", xerrors.Errorf("jwkutil: Generate use `%s` algorithm `%s` kid `%s`, unexpected key generation", use, algorithm, kid)
	}

	privBytes, err := privJWK.MarshalJSON()
	if err != nil {
		return "", "", xerrors.Errorf("jwkutil: Generate use `%s` algorithm `%s` marshal privJWK: %w", use, algorithm, err)
	}
	private = string(privBytes)

	pubBytes, err := pubJWK.MarshalJSON()
	if err != nil {
		return "", "", xerrors.Errorf("jwkutil: Generate use `%s` algorithm `%s` marshal pubJWK: %w", use, algorithm, err)
	}
	public = string(pubBytes)
	return
}
