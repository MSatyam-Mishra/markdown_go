package converter

import (
	"context"
	"io"
)

// Converter is the interface that wraps the Convert method for specific implementations.
type Converter interface {
	// Convert reads from r and returns the markdown representation.
	// It can also take an optional context and options.
	Convert(ctx context.Context, r io.Reader, opts *Options) (string, error)
}

// ConverterAPI is the core orchestrator interface that ZipConverter can use to recursively convert embedded files.
type ConverterAPI interface {
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
