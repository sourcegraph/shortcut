package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	makeRedirectURL, err := parseURLPattern("https://example.com/search?q=$QUERY&x=y")
	if err != nil {
		t.Fatal(err)
	}

	handler := handler{
		makeRedirectURL: makeRedirectURL,
		docsURL:         "https://example.com/docs",
	}

	tests := []struct {
		method       string
		url          string
		wantCode     int
		wantLocation string
	}{
		{
			method:       "GET",
			url:          "/foo bar",
			wantCode:     http.StatusFound,
			wantLocation: "https://example.com/search?q=foo+bar&x=y",
		},
		{
			method:       "HEAD",
			url:          "/foo bar",
			wantCode:     http.StatusFound,
			wantLocation: "https://example.com/search?q=foo+bar&x=y",
		},
		{
			method:       "GET",
			url:          "/repo:foo bar(",
			wantCode:     http.StatusFound,
			wantLocation: "https://example.com/search?q=repo%3Afoo+bar%28&x=y",
		},
		{
			method:       "GET",
			url:          "/",
			wantCode:     http.StatusFound,
			wantLocation: "https://example.com/docs",
		},
		{
			method:       "HEAD",
			url:          "/",
			wantCode:     http.StatusFound,
			wantLocation: "https://example.com/docs",
		},
		{
			method:   "POST",
			url:      "/foo bar",
			wantCode: http.StatusMethodNotAllowed,
		},
	}
	for _, test := range tests {
		t.Run(test.method+" "+test.url, func(t *testing.T) {
			req, err := http.NewRequest(test.method, test.url, nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if rr.Code != test.wantCode {
				t.Errorf("got code %d, want %d", rr.Code, test.wantCode)
			}
			if loc := rr.Header().Get("Location"); loc != test.wantLocation {
				t.Errorf("got location %q, want %q", loc, test.wantLocation)
			}
		})
	}
}
