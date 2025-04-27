package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"dream.website/internal/assert"
)

func TestSecureHeader(t *testing.T) {

	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	SecureHeader(next).ServeHTTP(rr, r)

	rs := rr.Result()

	expectedValue := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, expectedValue, rs.Header.Get("Content-Security-policy"))

	expectedValue = "origin-when-cross-origin"
	assert.Equal(t, expectedValue, rs.Header.Get("Referrer-Policy"))

	expectedValue = "nosniff"
	assert.Equal(t, expectedValue, rs.Header.Get("X-Content-Type-Options"))

	expectedValue = "deny"
	assert.Equal(t, rs.Header.Get("X-Frame-Options"), expectedValue)

	expectedValue = "0"
	assert.Equal(t, rs.Header.Get("X-XSS-Protection"), expectedValue)

	assert.Equal(t, http.StatusOK, rs.StatusCode)

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	assert.Equal(t, "OK", string(body))
}
