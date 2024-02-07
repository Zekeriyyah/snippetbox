package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zekeriyyah/snippetbox/internal/assert"
)

func TestSecureHeaders(t *testing.T) {
	rr := httptest.NewRecorder()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	secureHeaders(next).ServeHTTP(rr, r)

	rs := rr.Result()

	tests := []struct {
		actual   string
		expected string
	}{
		{
			actual:   "Content-Security-Policy",
			expected: "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com",
		},
		{
			actual:   "Referrer-Policy",
			expected: "origin-when-cross-origin",
		},
		{
			actual:   "X-Content-Options",
			expected: "nosniff",
		},
		{
			actual:   "X-Frame-Options",
			expected: "deny",
		},
		{
			actual:   "X-XSS-Protection",
			expected: "0",
		},
	}

	for _, testCase := range tests {
		assert.Equal(t, rs.Header.Get(testCase.actual), testCase.expected)
	}

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")

}
