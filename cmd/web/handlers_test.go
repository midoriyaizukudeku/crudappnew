package main

import (
	"net/http"
	"testing"

	"dream.website/internal/assert"
)

func TestPing(t *testing.T) {
	app := NewTestApplication(t)

	ts := newTestServer(t, app.Routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, "OK", body)
}

func TestSnippetView(t *testing.T) {

	app := NewTestApplication(t)

	ts := newTestServer(t, app.Routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid Id",
			urlPath:  "/snippet/view/1",
			wantCode: http.StatusOK,
			wantBody: "Of A Begiener",
		}, {
			name:     "Non Existent Id",
			urlPath:  "/snippet/view/2",
			wantCode: http.StatusNotFound,
		}, {
			name:     "decimal id",
			urlPath:  "/snippet/view/1.23",
			wantCode: http.StatusNotFound,
		}, {
			name:     "string id",
			urlPath:  "/snippet/view/foo",
			wantCode: http.StatusNotFound,
		}, {
			name:     "empty id",
			urlPath:  "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}

	for _, et := range tests {
		t.Run(et.name, func(t *testing.T) {
			code, _, body := ts.get(t, et.urlPath)
			t.Logf("Response body: %s", body) // Log the actual response body
			assert.Equal(t, et.wantCode, code)

			if et.wantBody != "" {
				// Check if the expected body is contained within the actual response
				assert.StringContains(t, body, et.wantBody)
			}
		})
	}

}

func TestUserSignu(t *testing.T) {
	app := NewTestApplication(t)
	ts := newTestServer(t, app.Routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/signup")
	t.Logf("Response body: %s", body) // Log the response body for debugging

	csrfTokenRx := extractcsrfToken(t, body)
	t.Logf("csrf token is %q", csrfTokenRx)

}
