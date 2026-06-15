package main

import (
	"net/http"
	"testing"

	"github.com/ridhamu/snippetbox/internal/assert"
)

func TestPingHandler(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	httpStatusCode, _, httpResponseBody := ts.get(t, "/ping")

	assert.Equal(t, httpStatusCode, http.StatusOK)
	assert.Equal(t, httpResponseBody, "OK")
}
