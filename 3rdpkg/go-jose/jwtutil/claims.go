package jwtutil

import (
	"encoding/json"

	"github.com/square/go-jose/v3/jwt"
	"golang.org/x/xerrors"
)

// Claims holds standard and custom claims
type Claims struct {
	Standard *jwt.Claims
	Custom   map[string]interface{}
}

// UnmarshalJSON implements json.Unmarshaler
func (c *Claims) UnmarshalJSON(b []byte) error {
	if b == nil {
		return nil
	}
	var entries map[string]interface{}
	if err := json.Unmarshal(b, &entries); err != nil {
		return xerrors.Errorf("jwtutil: (*Claims).UnmarshalJSON: %w", err)
	}

	if c.Standard == nil {
		c.Standard = &jwt.Claims{}
	}
	if c.Custom == nil {
		c.Custom = make(map[string]interface{})
	}

	for k, v := range entries {
		switch k {
		case "iss":
			v, ok := v.(string)
			if !ok {
				return xerrors.New("jwtutil: (*Claims).UnmarshalJSON iss must be a string")
			}
			c.Standard.Issuer = v
		case "sub":
			v, ok := v.(string)
			if !ok {
				return xerrors.New("jwtutil: (*Claims).UnmarshalJSON sub must be a string")
			}
			c.Standard.Subject = v
		case "aud":
			switch v := v.(type) {
			case string:
				aud := jwt.Audience{v}
				c.Standard.Audience = aud
			case []string:
				c.Standard.Audience = jwt.Audience(v)
			case []interface{}:
				aud := make(jwt.Audience, len(v))
				for i, e := range v {
					e, ok := e.(string)
					if !ok {
						return xerrors.Errorf("jwtutil: (*Claims).UnmarshalJSON unexpected type of aud entry %T", e)
					}
					aud[i] = e
				}
				c.Standard.Audience = aud
			default:
				return xerrors.Errorf("jwtutil: (*Claims).UnmarshalJSON unexpected type of aud %T", v)
			}
		case "exp":
			switch v := v.(type) {
			case int64:
				exp := jwt.NumericDate(v)
				c.Standard.Expiry = &exp
			case float64:
				exp := jwt.NumericDate(v)
				c.Standard.Expiry = &exp
			default:
				return xerrors.Errorf("jwtutil: (*Claims).UnmarshalJSON unexpected type of exp %T", v)
			}
		case "nbf":
			switch v := v.(type) {
			case int64:
				nbf := jwt.NumericDate(v)
				c.Standard.NotBefore = &nbf
			case float64:
				nbf := jwt.NumericDate(v)
				c.Standard.NotBefore = &nbf
			default:
				return xerrors.Errorf("jwtutil: (*Claims).UnmarshalJSON unexpected type of nbf %T", v)
			}
		case "iat":
			switch v := v.(type) {
			case int64:
				iat := jwt.NumericDate(v)
				c.Standard.IssuedAt = &iat
			case float64:
				iat := jwt.NumericDate(v)
				c.Standard.IssuedAt = &iat
			default:
				return xerrors.Errorf("jwtutil: (*Claims).UnmarshalJSON unexpected type of iat %T", v)
			}
		case "jti":
			v, ok := v.(string)
			if !ok {
				return xerrors.New("jwtutil: (*Claims).UnmarshalJSON jti must be a string")
			}
			c.Standard.ID = v
		default:
			c.Custom[k] = v
		}
	}
	return nil
}

// MarshalJSON implements json.marshaler
func (c Claims) MarshalJSON() ([]byte, error) {
	const standardClaimFields = 7
	entries := make(map[string]interface{}, len(c.Custom)+standardClaimFields)

	for k, v := range c.Custom {
		entries[k] = v
	}
	if c.Standard.Issuer != "" {
		entries["iss"] = c.Standard.Issuer
	}
	if c.Standard.Subject != "" {
		entries["sub"] = c.Standard.Subject
	}
	if len(c.Standard.Audience) > 0 {
		entries["aud"] = c.Standard.Audience
	}
	if c.Standard.Expiry != nil {
		entries["exp"] = c.Standard.Expiry
	}
	if c.Standard.NotBefore != nil {
		entries["nbf"] = c.Standard.NotBefore
	}
	if c.Standard.IssuedAt != nil {
		entries["iat"] = c.Standard.IssuedAt
	}
	if c.Standard.ID != "" {
		entries["jti"] = c.Standard.ID
	}

	b, err := json.Marshal(entries)
	if err != nil {
		return nil, xerrors.Errorf("jwtutil: (Claims).MarshalJSON: %w", err)
	}
	return b, nil
}
