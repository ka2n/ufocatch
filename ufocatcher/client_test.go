package ufocatcher

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var mockHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != "/atom/edinetx/query/4751" {
		http.Error(w, fmt.Sprintf("invalid request URI: %v", r.RequestURI), 400)
		return
	}
	f, err := os.Open("./test_file/edinetx.xml")
	if err != nil {
		panic(err)
	}
	io.Copy(w, f)
})

func TestGet(t *testing.T) {
	ts := httptest.NewServer(mockHandler)
	defer ts.Close()

	ctx := context.Background()
	feed, err := Get(ctx, ts.URL, CategoryEdinetx, "4751")
	if err != nil {
		t.Fatal(err)
	}
	if len(feed.Entries) == 0 {
		t.Errorf("expected at least 1 entry, but: %v", len(feed.Entries))
	}
}
