package main

import (
	"bytes"
	"html"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"dream.website/internal/model/mocks"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
)

var csrfTokenRx = regexp.MustCompile(`<input type="hidden" name="csrf_token" value="(.+?)">`)

func extractcsrfToken(t *testing.T, body string) string {
	matches := csrfTokenRx.FindStringSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}
	return html.UnescapeString(string(matches[1]))
}

func NewTestApplication(t *testing.T) *Application {

	templateCache, err := NewTemplatecache()
	if err != nil {
		t.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	sessionManeger := scs.New()
	sessionManeger.Lifetime = 12 * time.Hour
	sessionManeger.Cookie.Secure = true

	return &Application{
		errorLog:       log.New(io.Discard, "", 0),
		infoLog:        log.New(io.Discard, "", 0),
		snippets:       &mocks.SnippetModel{},
		Users:          &mocks.UserModel{},
		Templatecache:  templateCache,
		FormDecoder:    formDecoder,
		SessionManager: sessionManeger,
	}

}

type TestServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *TestServer {
	ts := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &TestServer{ts}
}

func (ts *TestServer) get(t *testing.T, urlPath string) (int, http.Header, string) {

	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)

}
