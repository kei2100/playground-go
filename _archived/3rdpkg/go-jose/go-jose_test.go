package jose_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"log"
	"testing"

	"github.com/square/go-jose/v3"
)

func TestFoo(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	pubJWK := jose.JSONWebKey{Key: pub, KeyID: "myid", Algorithm: "EdDSA", Use: "sig"}
	privJWK := jose.JSONWebKey{Key: priv, KeyID: "myid", Algorithm: "EdDSA", Use: "sig"}

	if privJWK.IsPublic() || !pubJWK.IsPublic() || !privJWK.Valid() || !pubJWK.Valid() {
		t.Fatal("omg")
	}

	pubJSON, err := pubJWK.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	privJSON, err := privJWK.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("\n=== pub ===\n%s\n=== priv ===\n%s", string(pubJSON), string(privJSON))

	var unmarshalPubJWK jose.JSONWebKey
	if err := unmarshalPubJWK.UnmarshalJSON(pubJSON); err != nil {
		t.Fatal(err)
	}
	var unmarshalPrivJWK jose.JSONWebKey
	if err := unmarshalPrivJWK.UnmarshalJSON(privJSON); err != nil {
		t.Fatal(err)
	}

	log.Printf("\n=== pub ===\n%+v\n=== priv ===\n%+v", unmarshalPubJWK, unmarshalPrivJWK)

	pubJSON2, _ := unmarshalPubJWK.MarshalJSON()
	privJSON2, _ := unmarshalPrivJWK.MarshalJSON()

	log.Printf("\n=== pub ===\n%s\n=== priv ===\n%s", string(pubJSON2), string(privJSON2))
}
