package jwtutil_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/kei2100/playground-go/3rdpkg/go-jose/jwtutil"

	"github.com/square/go-jose/v3/jwt"
)

func numericDatep(i int) *jwt.NumericDate {
	v := jwt.NumericDate(i)
	return &v
}

func TestClaims_UnmarshalJSON(t *testing.T) {

	tt := []struct {
		json       string
		wantClaims jwtutil.Claims
	}{
		{
			json: `
					{
						"iss": "my-issuer",
						"sub": "my-subject",
						"aud": "my-audience",
						"exp": 1571462105,
						"nbf": 1571462104,
						"iat": 1571462103,
						"jti": "abcd",
						"custom": "my-custom"
					}
			`,
			wantClaims: jwtutil.Claims{
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
			},
		},
		{
			json: `{"aud": ["foo","bar"]}`,
			wantClaims: jwtutil.Claims{
				Standard: &jwt.Claims{
					Audience: jwt.Audience{"foo", "bar"},
				},
			},
		},
	}
	for i, te := range tt {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			var gotClaims jwtutil.Claims
			if err := json.Unmarshal([]byte(te.json), &gotClaims); err != nil {
				t.Error(err)
			}
			if te.wantClaims.Standard != nil {
				got, want := gotClaims.Standard, te.wantClaims.Standard
				if g, w := got.Issuer, want.Issuer; g != w {
					t.Errorf("Issuer got %v, want %v", g, w)
				}
				if g, w := got.Subject, want.Subject; g != w {
					t.Errorf("Subject got %v, want %v", g, w)
				}
				if !reflect.DeepEqual(got.Audience, want.Audience) {
					t.Errorf("Audience got %v, want %v", got.Audience, want.Audience)
				}
				if want.Expiry != nil {
					if g, w := *got.Expiry, *want.Expiry; g != w {
						t.Errorf("Expiry got %v, want %v", g, w)
					}
				}
				if want.NotBefore != nil {
					if g, w := *got.NotBefore, *want.NotBefore; g != w {
						t.Errorf("NotBefore got %v, want %v", g, w)
					}
				}
				if want.IssuedAt != nil {
					if g, w := *got.IssuedAt, *want.IssuedAt; g != w {
						t.Errorf("IssuedAt got %v, want %v", g, w)
					}
				}
				if g, w := got.ID, want.ID; g != w {
					t.Errorf("ID got %v, want %v", g, w)
				}
			}
			if len(te.wantClaims.Custom) > 0 {
				if !reflect.DeepEqual(gotClaims.Custom, te.wantClaims.Custom) {
					t.Errorf("Custom got\n%+v\nwant\n%+v\n", gotClaims.Custom, te.wantClaims.Custom)
				}
			}
		})
	}
}

func TestClaims_MarshalJSON(t *testing.T) {
	tt := []struct {
		claims   jwtutil.Claims
		wantJSON string
	}{
		{
			claims: jwtutil.Claims{
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
			},
			wantJSON: `
					{
						"iss": "my-issuer",
						"sub": "my-subject",
						"aud": "my-audience",
						"exp": 1571462105,
						"nbf": 1571462104,
						"iat": 1571462103,
						"jti": "abcd",
						"custom": "my-custom"
					}
			`,
		},
		{
			claims: jwtutil.Claims{
				Standard: &jwt.Claims{
					Audience: jwt.Audience{"foo", "bar"},
				},
			},
			wantJSON: `{"aud": ["foo","bar"]}`,
		},
	}
	for i, te := range tt {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			b, err := json.Marshal(te.claims)
			if err != nil {
				t.Fatal(err)
			}
			var gotClaims map[string]interface{}
			if err := json.Unmarshal(b, &gotClaims); err != nil {
				t.Fatal(err)
			}
			var wantClaims map[string]interface{}
			if err := json.Unmarshal([]byte(te.wantJSON), &wantClaims); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(gotClaims, wantClaims) {
				t.Errorf("Claims got\n%+v\nwant\n%+v\n", gotClaims, wantClaims)
			}
		})
	}

}
