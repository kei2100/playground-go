package go1_21

import (
	"log/slog"
	"net/url"
	"os"
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
	return slog.GroupValue(
		slog.String("token", "xxxxxx"),
		slog.Any("data", v.Data),
	)
}

func (v *CredDataLogValuer) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("id", v.ID),
		slog.String("password", "yyyyyy"),
	)
}
