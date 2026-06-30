package markdown_go

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// MarkItDown is the main orchestrator for conversions.
type MarkItDown struct {
	converters map[string]Converter // map of extension to converter
}

// NewMarkItDown initializes a new MarkItDown instance with default converters.
func NewMarkItDown() *MarkItDown {
	m := &MarkItDown{
		converters: make(map[string]Converter),
	}
	m.registerDefaultConverters()
	return m
}

func (m *MarkItDown) registerDefaultConverters() {
	// HTML
	htmlConv := &HTMLConverter{}
	m.converters[".html"] = htmlConv
	m.converters[".htm"] = htmlConv

	// Office and PDF
	docConv := &DocConverter{}
	m.converters[".pdf"] = docConv
	m.converters[".docx"] = docConv
	m.converters[".pptx"] = docConv
	m.converters[".pages"] = docConv

	// Data
	dataConv := &DataConverter{}
	m.converters[".csv"] = dataConv
	m.converters[".json"] = dataConv
	m.converters[".xml"] = dataConv

	// Zip
	m.converters[".zip"] = &ZipConverter{markItDown: m}

	// Images
	imgConv := &ImageConverter{}
	m.converters[".jpg"] = imgConv
	m.converters[".jpeg"] = imgConv
	m.converters[".png"] = imgConv
	m.converters[".gif"] = imgConv

	// Audio
	audioConv := &AudioConverter{}
	m.converters[".mp3"] = audioConv
	m.converters[".wav"] = audioConv
	m.converters[".ogg"] = audioConv
	m.converters[".flac"] = audioConv

	// EPub
	epubConv := &EpubConverter{htmlConverter: htmlConv}
	m.converters[".epub"] = epubConv
}

// ConvertFile reads a file and converts it to markdown based on its extension.
func (m *MarkItDown) ConvertFile(ctx context.Context, path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(path))
	opts := &Options{
		Extension: ext,
		FileName:  filepath.Base(path),
	}

	return m.Convert(ctx, file, opts)
}

// ConvertURL handles URL inputs, delegating YouTube URLs to YoutubeConverter.
func (m *MarkItDown) ConvertURL(ctx context.Context, urlStr string) (string, error) {
	if strings.Contains(urlStr, "youtube.com") || strings.Contains(urlStr, "youtu.be") {
		ytConv := &YoutubeConverter{}
		opts := &Options{URL: urlStr}
		return ytConv.Convert(ctx, nil, opts)
	}
	// For other URLs, you would perform an HTTP GET, detect the MIME type, 
	// map it to an extension, and call m.Convert.
	return "", fmt.Errorf("unsupported URL or non-youtube URL not implemented yet")
}

// Convert processes an io.Reader stream using the appropriate converter based on Options.
func (m *MarkItDown) Convert(ctx context.Context, r io.Reader, opts *Options) (string, error) {
	if opts == nil {
		opts = &Options{}
	}

	conv, ok := m.converters[opts.Extension]
	if !ok {
		return "", fmt.Errorf("unsupported file extension: %s", opts.Extension)
	}

	return conv.Convert(ctx, r, opts)
}
