package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/ridhamu/snippetbox/internal/models/mocks"
)

func newTestApplication(t *testing.T) *application {
	templateCacheTest, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	formDecoderTest := form.NewDecoder()

	sessionManagerTest := scs.New()
	sessionManagerTest.Lifetime = 12 * time.Hour
	sessionManagerTest.Cookie.Secure = true

	return &application{
		logger:         slog.New(slog.NewTextHandler(io.Discard, nil)),
		snippetModel:   &mocks.SnippetModel{},
		userModel:      &mocks.UserModel{},
		templateCache:  templateCacheTest,
		formDecoder:    formDecoderTest,
		sessionManager: sessionManagerTest,
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	// add cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	result, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatal(err)
	}

	body = bytes.TrimSpace(body)

	return result.StatusCode, result.Header, string(body)
}
