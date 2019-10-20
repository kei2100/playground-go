package jwtutil_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/kei2100/playground-go/3rdpkg/go-jose/jwkutil"
	"github.com/kei2100/playground-go/3rdpkg/go-jose/jwtutil"

	"github.com/rs/xid"
	"github.com/square/go-jose/v3/jwt"
)

func TestNewSigned_Verify(t *testing.T) {
	kid1 := xid.New().String()
	pub1, priv1, err := jwkutil.Generate("sig", "EdDSA", kid1)
	if err != nil {
		t.Fatal(err)
	}
	kid2 := xid.New().String()
	pub2, priv2, err := jwkutil.Generate("sig", "EdDSA", kid2)
	if err != nil {
		t.Fatal(err)
	}
	kidSet1_2 := fmt.Sprintf(`{"keys":[%s,%s]}`, pub1, pub2)

	kid3 := xid.New().String()
	pub3, priv3, err := jwkutil.Generate("sig", "EdDSA", kid3)
	if err != nil {
		t.Fatal(err)
	}
	kidSet3 := fmt.Sprintf(`{"keys":[%s]}`, pub3)

	jwt1, err := jwtutil.NewSigned(priv1, jwtutil.Claims{Standard: &jwt.Claims{Subject: "test"}})
	if err != nil {
		t.Fatal(err)
	}
	jwt2, err := jwtutil.NewSigned(priv2, jwtutil.Claims{Standard: &jwt.Claims{Subject: "test"}})
	if err != nil {
		t.Fatal(err)
	}
	jwt3, err := jwtutil.NewSigned(priv3, jwtutil.Claims{Standard: &jwt.Claims{Subject: "test"}})
	if err != nil {
		t.Fatal(err)
	}

	tt := []struct {
		jwtString    string
		jwkSetString string
		wantOK       bool
	}{
		{
			jwtString:    jwt1,
			jwkSetString: kidSet1_2,
			wantOK:       true,
		},
		{
			jwtString:    jwt2,
			jwkSetString: kidSet1_2,
			wantOK:       true,
		},
		{
			jwtString:    jwt3,
			jwkSetString: kidSet1_2,
			wantOK:       false,
		},
		{
			jwtString:    jwt1,
			jwkSetString: kidSet3,
			wantOK:       false,
		},
		{
			jwtString:    jwt2,
			jwkSetString: kidSet3,
			wantOK:       false,
		},
		{
			jwtString:    jwt3,
			jwkSetString: kidSet3,
			wantOK:       true,
		},
	}
	for i, te := range tt {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			var claims jwtutil.Claims
			err := jwtutil.Verify(te.jwtString, te.jwkSetString, &claims)
			if te.wantOK && err != nil {
				t.Errorf("got %v, want no error", err)
			}
			if !te.wantOK && err == nil {
				t.Error("got no error, want an error")
			}
		})
	}
}

func TestUnsafeClaims(t *testing.T) {
	_, priv, err := jwkutil.Generate("sig", "EdDSA", xid.New().String())
	if err != nil {
		t.Fatal(err)
	}
	jwtString, err := jwtutil.NewSigned(priv, jwtutil.Claims{
		Standard: &jwt.Claims{
			Issuer:    "my-issuer",
			Subject:   "my-subject",
			Audience:  jwt.Audience{"my-audience"},
			Expiry:    numericDatep(1571462105),
			NotBefore: numericDatep(1571462104),
			IssuedAt:  numericDatep(1571462103),
			ID:        "abcd",
		},
		Custom: map[string]interface{}{
			"custom": "my-custom",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	var got jwtutil.Claims
	if err := jwtutil.UnsafeClaims(jwtString, &got); err != nil {
		t.Fatal(err)
	}
	if g, w := got.Standard.Issuer, "my-issuer"; g != w {
		t.Errorf("Issuer got %v, want %v", g, w)
	}
	if g, w := got.Standard.Subject, "my-subject"; g != w {
		t.Errorf("Subject got %v, want %v", g, w)
	}
	if !reflect.DeepEqual(got.Standard.Audience, jwt.Audience{"my-audience"}) {
		t.Errorf("Audience got %v, want %v", got.Standard.Audience, jwt.Audience{"my-audience"})
	}
	if g, w := got.Standard.Expiry.Time().Unix(), int64(1571462105); g != w {
		t.Errorf("Expiry got %v, want %v", g, w)
	}
	if g, w := got.Standard.NotBefore.Time().Unix(), int64(1571462104); g != w {
		t.Errorf("NotBefore got %v, want %v", g, w)
	}
	if g, w := got.Standard.IssuedAt.Time().Unix(), int64(1571462103); g != w {
		t.Errorf("IssuedAt got %v, want %v", g, w)
	}
	if g, w := got.Standard.ID, "abcd"; g != w {
		t.Errorf("ID got %v, want %v", g, w)
	}
	if len(got.Custom) == 0 {
		t.Error("got.Custom length is zero")
	} else {
		if g, w := got.Custom["custom"].(string), "my-custom"; g != w {
			t.Errorf("got.Custom[`custom`] got %v, want %v", g, w)
		}
	}
}
