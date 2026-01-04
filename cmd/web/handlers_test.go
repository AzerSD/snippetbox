package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"snippetbox.azersd.me/internal/assert"
	"testing"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")
}