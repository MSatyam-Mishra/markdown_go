package markdown_go

import (
	"context"
	"fmt"
	"io"
	"net/http"
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

// ConvertURL handles URL inputs, delegating YouTube URLs to YoutubeConverter, or fetching general webpages.
func (m *MarkItDown) ConvertURL(ctx context.Context, urlStr string) (string, error) {
	if strings.Contains(urlStr, "youtube.com") || strings.Contains(urlStr, "youtu.be") {
		ytConv := &YoutubeConverter{}
		opts := &Options{URL: urlStr}
		return ytConv.Convert(ctx, nil, opts)
	}

	// Fetch generic web page
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	// If it's HTML, convert it
	if strings.Contains(strings.ToLower(contentType), "text/html") {
		htmlConv := &HTMLConverter{}
		return htmlConv.Convert(ctx, resp.Body, &Options{Extension: ".html", URL: urlStr})
	}

	return "", fmt.Errorf("unsupported content type from URL: %s", contentType)
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
