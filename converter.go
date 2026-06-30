package markdown_go

import (
	"context"
	"io"
)

// Converter is the interface that wraps the Convert method.
type Converter interface {
	// Convert reads from r and returns the markdown representation.
	// It can also take an optional context and options.
	Convert(ctx context.Context, r io.Reader, opts *Options) (string, error)
}

// Options provides configuration for the conversion process.
type Options struct {
	// Extension is the original file extension (e.g., ".pdf", ".html")
	Extension string
	// FileName is the original filename, if available.
	FileName string
	// URL is the source URL, if applicable.
	URL string
}
