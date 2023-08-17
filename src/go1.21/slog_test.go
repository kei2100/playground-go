package go1_21

import (
	"context"
	"log/slog"
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestSlog_Simple(t *testing.T) {
	slog.Info("message", "foo", "bar")
	// 2023/07/08 15:39:13 INFO message foo=bar
}

func TestSlog_JSONMessage(t *testing.T) {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "password" {
				a.Value = slog.StringValue("******")
			}
			return a
		},
	}))
	log.Debug("logged in", "id", "id:foo", "password", "sensitive value")
	// {"time":"2023-07-08T15:45:41.097768+09:00","level":"DEBUG","msg":"logged in","id":"id:foo","password":"******"}

	// To set as default logger, do the following
	//slog.SetDefault(log)
}

func TestSlog_With(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	foosLog := log.With("id", "id:foo")
	barsLog := log.With("id", "id:bar")

	foosLog.Info("logged out")
	barsLog.Info("logged out")
	// time=2023-07-08T15:51:10.253+09:00 level=INFO msg="logged out" id=id:foo
	// time=2023-07-08T15:51:10.253+09:00 level=INFO msg="logged out" id=id:bar
}

func TestSlog_Group(t *testing.T) {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	log.Info("send request", slog.Group(
		"request",
		"method", "GET",
		"path", "/foo/bar",
	))
	// {"time":"2023-07-08T15:58:24.69949+09:00","level":"INFO","msg":"send request","request":{"method":"GET","path":"/foo/bar"}}
}

func TestSlog_PerformanceConsiderations(t *testing.T) {
	u, _ := url.Parse("https://example.com/foo?bar=baz")

	slog.Debug("starting request", "url", u.String()) // BAD: may compute u.String() unnecessarily
	slog.Debug("starting request", "url", u)          // GOOD: calls u.String() only if needed
}

func TestSlog_PrintStruct(t *testing.T) {
	c := &Cred{
		Token: "sensitive",
		Data: &CredData{
			ID:       "foo",
			Password: "password",
		},
	}
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	log.Info("message", "cred", c)
	// {"time":"2023-07-08T16:16:15.814178+09:00","level":"INFO","msg":"message","cred":{"Token":"sensitive","Data":{"ID":"foo","Password":"password"}}}

	c2 := &Cred{
		Token: "sensitive",
		Data:  nil,
	}
	log.Info("message", "cred2", c2)
	// {"time":"2023-07-08T17:16:02.578957+09:00","level":"INFO","msg":"message","cred2":{"Token":"sensitive","Data":null}}
}

func TestSlog_Value(t *testing.T) {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	log.Info("test", "any_nil", slog.AnyValue(nil))
	// {"time":"2023-07-08T17:21:34.1266+09:00","level":"INFO","msg":"test","any_nil":null}
}

type Cred struct {
	Token string
	Data  *CredData
}

type CredData struct {
	ID       string
	Password string
}

func TestSlog_PrintLogValuer(t *testing.T) {
	c := &CredLogValuer{
		Token: "sensitive",
		Data: &CredDataLogValuer{
			ID:       "foo",
			Password: "password",
		},
	}
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	log.Info("message", "cred", c)
	// {"time":"2023-07-08T16:27:58.235297+09:00","level":"INFO","msg":"message","cred":{"token":"xxxxxx","data":{"id":"foo","password":"yyyyyy"}}}
	c = &CredLogValuer{
		Token: "sensitive",
		Data:  nil,
	}
	log.Info("message", "cred", c)
	// {"time":"2023-07-19T20:15:10.945092+09:00","level":"INFO","msg":"message","cred":{"token":"xxxxxx","data":null}}
}

type CredLogValuer struct {
	Token string
	Data  *CredDataLogValuer
}

type CredDataLogValuer struct {
	ID       string
	Password string
}

func (v *CredLogValuer) LogValue() slog.Value {
	if v == nil {
		return slog.AnyValue(nil)
	}
	return slog.GroupValue(
		slog.String("token", "xxxxxx"),
		slog.Attr{Key: "data", Value: v.Data.LogValue()},
	)
}

func (v *CredDataLogValuer) LogValue() slog.Value {
	if v == nil {
		return slog.AnyValue(nil)
	}
	slog.Value{}.Kind()
	return slog.GroupValue(
		slog.String("id", v.ID),
		slog.String("password", "yyyyyy"),
	)
}

func TestSlog_PrintCollection(t *testing.T) {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	log.LogAttrs(
		context.Background(),
		slog.LevelInfo,
		"message",
		slog.String("scalar", "value"),
		slog.Any("empty_slice", make([]string, 0)),
		slog.Any("slice", []float32{
			0.1,
			0.2,
		}),
		slog.Any("empty_map", make(map[string]string)),
		slog.Any("empty_map", map[string]float32{
			"one": 0.1,
			"two": 0.2,
		}),
		slog.Any("int_map", map[int]float32{
			10: 0.1,
			20: 0.2,
		}),
		slog.Any("log_valuer_slice", []string{
			upperCaseLogValuer("aaa").LogValue().String(),
			alwaysZeroLogValur(100).LogValue().String(),
		}),
	)
	// {"time":"2023-08-17T15:52:55.558225+09:00","level":"INFO","msg":"message","scalar":"value","empty_slice":[],"slice":[0.1,0.2],"empty_map":{},"empty_map":{"one":0.1,"two":0.2},"int_map":{"10":0.1,"20":0.2},"log_valuer_slice":["AAA","0"]}
}

type upperCaseLogValuer string

func (v upperCaseLogValuer) LogValue() slog.Value {
	return slog.StringValue(strings.ToUpper(string(v)))
}

type alwaysZeroLogValur int64

func (v alwaysZeroLogValur) LogValue() slog.Value {
	return slog.Int64Value(0)
}
