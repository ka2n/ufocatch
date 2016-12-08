package ufocatch

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var mockHandler = http.NewServeMux()

func init() {
	serveFile := func(fname string, w http.ResponseWriter) {
		f, err := os.Open(fname)
		if err != nil {
			panic(err)
		}
		io.Copy(w, f)
	}
	mockHandler.HandleFunc("/atom/edinetx/query/4751", func(w http.ResponseWriter, r *http.Request) {
		serveFile("./test_file/edinetx.xml", w)
	})
	mockHandler.HandleFunc("/atom/edinetx", func(w http.ResponseWriter, r *http.Request) {
		serveFile("./test_file/edinetx.xml", w)
	})
}

func TestGet(t *testing.T) {
	ts := httptest.NewServer(mockHandler)
	defer ts.Close()

	ctx := context.Background()
	client, err := New(ts.URL, http.DefaultClient)
	if err != nil {
		t.Fatal(err)
	}
	feed, err := client.Get(ctx, CategoryEdinetx, "4751")
	if err != nil {
		t.Fatal(err)
	}
	if len(feed.Entries) == 0 {
		t.Errorf("expected at least 1 entry, but: %v", len(feed.Entries))
	}

	feed, err = client.Get(ctx, CategoryEdinetx, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(feed.Entries) == 0 {
		t.Errorf("expected at least 1 entry, but: %v", len(feed.Entries))
	}
}
