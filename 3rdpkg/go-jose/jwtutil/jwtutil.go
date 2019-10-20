package jwtutil

import (
	"encoding/json"

	"github.com/square/go-jose/v3"
	"github.com/square/go-jose/v3/jwt"
	"golang.org/x/xerrors"
)

// NewSigned generates a new signed JWT that signed by specified private jwk and claims
func NewSigned(jwkString string, claims Claims) (string, error) {
	var jwk jose.JSONWebKey
	if err := jwk.UnmarshalJSON([]byte(jwkString)); err != nil {
		return "", xerrors.Errorf("jwtutil: NewSigned unmarshal jwk: %w", err)
	}
	if !jwk.Valid() || jwk.IsPublic() {
		return "", xerrors.Errorf("jwtutil: NewSigned invalid jwk")
	}
	if !allowSignAlgorithm(jwk.Algorithm) {
		return "", xerrors.Errorf("jwtutil: NewSigned alg %s is not allowed", jwk.Algorithm)
	}

	var opt jose.SignerOptions
	opt.WithType("JWT")
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.SignatureAlgorithm(jwk.Algorithm), Key: jwk}, &opt)
	if err != nil {
		return "", xerrors.Errorf("jwtutil: NewSigned create signer: %w", err)
	}

	builder := jwt.Signed(signer)
	if len(claims.Custom) > 0 {
		builder = builder.Claims(claims.Custom)
	}
	if claims.Standard != nil {
		builder = builder.Claims(claims.Standard)
	}
	jwtString, err := builder.CompactSerialize()
	if err != nil {
		return "", xerrors.Errorf("jwtutil: NewSigned serialize: %w", err)
	}
	return jwtString, nil
}

// Verify and de-serialize a JWT into dest using the provided JWK-Set
func Verify(jwtString, jwkSetString string, dest *Claims) error {
	var jwkSet jose.JSONWebKeySet
	if err := json.Unmarshal([]byte(jwkSetString), &jwkSet); err != nil {
		return xerrors.Errorf("jwtutil: Verify unmarshal jwkSet: %w", err)
	}
	for _, jwk := range jwkSet.Keys {
		if !jwk.Valid() || !jwk.IsPublic() {
			return xerrors.Errorf("jwtutil: Verify invalid jwk")
		}
		if !allowSignAlgorithm(jwk.Algorithm) {
			return xerrors.Errorf("jwtutil: Verify alg %s is not allowed", jwk.Algorithm)
		}
	}

	token, err := jwt.ParseSigned(jwtString)
	if err != nil {
		return xerrors.Errorf("jwtutil: Verify parse jwt: %w", err)
	}
	if err := token.Claims(&jwkSet, dest); err != nil {
		return xerrors.Errorf("jwtutil: Verify parse and verify claims: %w", err)
	}
	return nil
}

// UnsafeClaims de-serializes the claims of a JWT into the dest.
// For signed JWTs, the claims are not verified.
// This function won't work for encrypted JWTs.
func UnsafeClaims(jwtString string, dest *Claims) error {
	token, err := jwt.ParseSigned(jwtString)
	if err != nil {
		return xerrors.Errorf("jwtutil: UnsafeClaims parse jwt: %w", err)
	}
	if err := token.UnsafeClaimsWithoutVerification(dest); err != nil {
		return xerrors.Errorf("jwtutil: UnsafeClaims unmarshal: %w", err)
	}
	return nil
}

func allowSignAlgorithm(algorithm string) bool {
	switch algorithm {
	case "EdDSA":
		return true
	default:
		return false
	}
}
