package ufocatch

import "context"

// Client retrieves Feed
type Client interface {
	Get(ctx context.Context, endpoint Endpoint, cat Category, query string) (*Feed, error)
	Download(ctx context.Context, endpoint Endpoint, format Format, id string) (string, error)
}

// Category for resource
type Category string

// List of categories
const (
	CategoryEdinet  Category = "edinet"
	CategoryEdinetx          = "edinetx"
	CategoryTdnet            = "tdnet"
	CategoryTdnetx           = "tdnetx"
	CategoryCg               = "cg"
)

// Endpoint for request
type Endpoint string

const (
	// DefaultEndpoint for ufocatch API
	DefaultEndpoint Endpoint = "http://resource.ufocatch.com"
)

// Format of file to download
type Format string

// List of formats
const (
	FormatData Format = "data"
	FormatPDF         = "pdf"
)
