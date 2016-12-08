package ufocatch

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// New Client
func New(endpoint string, client *http.Client) (*Client, error) {
	if endpoint == "" {
		endpoint = DefaultEndpoint
	}

	if client == nil {
		client = http.DefaultClient
	}

	return &Client{
		Endpoint: endpoint,
		Client:   client,
	}, nil
}

// Client retrieves Feed from ufocatcher
type Client struct {
	Endpoint string
	Client   *http.Client
}

// Get /atom/{種別}/query/{クエリワード}
func (c Client) Get(ctx context.Context, cat string, query string) (*Feed, error) {
	p := path.Join("/atom", cat)
	if query != "" {
		p = path.Join(p, "query", query)
	}
	req, err := http.NewRequest("GET", c.Endpoint+p, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	r, err := c.Client.Do(req)

	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(r.Body)
		return nil, fmt.Errorf("invalid response: code: %v, body: %v", r.StatusCode, string(b))
	}

	var feed Feed
	if err := xml.NewDecoder(r.Body).Decode(&feed); err != nil {
		return nil, err
	}

	return &feed, nil
}

// Download /{format: pdf, data}/{source: edinet/tdnet}/{id}
func (c Client) Download(ctx context.Context, format string, id string) (string, error) {
	source := sourceByID(id)
	if source == "" {
		return "", fmt.Errorf("unknown format for id: '%v'", id)
	}

	// Create request
	urlStr := c.Endpoint + path.Join("/", format, source, id)
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", err
	}
	req = req.WithContext(ctx)

	// Execute request
	r, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		if r.StatusCode == http.StatusNotFound {
			return "", fmt.Errorf("404 file not found: %v", urlStr)
		}
		b, _ := ioutil.ReadAll(r.Body)
		return "", fmt.Errorf("invalid response: code: %v, url: %v, body: %v", r.StatusCode, urlStr, string(b))
	}

	// Determine filename
	fileName := fileNameByHeader(r.Header, id)

	// Create file
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		os.Remove(fileName)
	}
	f, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer f.Close()
	io.Copy(f, r.Body)

	return fileName, nil
}

func fileNameByHeader(h http.Header, base string) string {
	if t, params, _ := mime.ParseMediaType(h.Get("Content-Disposition")); t == "attachment" {
		if ext := filepath.Ext(params["filename"]); ext != "" {
			return base + ext
		}
	}
	if t, _, _ := mime.ParseMediaType(h.Get("Content-Type")); t != "" {
		switch {
		case strings.Contains(t, "pdf"):
			return base + ".pdf"
		case strings.Contains(t, "zip"):
			return base + ".zip"
		}
	}
	return base
}

func sourceByID(id string) string {
	pre := id[0:2]
	switch pre {
	case "TD":
		return "tdnet"
	case "ED":
		return "edinet"
	case "CG":
		return "cg"
	}
	return ""
}
