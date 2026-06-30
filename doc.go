package markdown_go

import (
	"bytes"
	"context"
	"io"
	"strings"

	"code.sajari.com/docconv"
	"github.com/ledongthuc/pdf"
)

type DocConverter struct{}

func (c *DocConverter) Convert(ctx context.Context, r io.Reader, opts *Options) (string, error) {
	if opts != nil && opts.Extension == ".pdf" {
		return c.convertPDF(r)
	}

	mimeType := ""
	if opts != nil && opts.Extension != "" {
		mimeType = docconv.MimeTypeByExtension(opts.Extension)
	}

	res, err := docconv.Convert(r, mimeType, true)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(res.Body), nil
}

func (c *DocConverter) convertPDF(r io.Reader) (string, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	bytesReader := bytes.NewReader(b)
	pdfReader, err := pdf.NewReader(bytesReader, int64(len(b)))
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	for i := 1; i <= pdfReader.NumPage(); i++ {
		page := pdfReader.Page(i)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			continue
		}

		sb.WriteString(text)
		sb.WriteString("\n\n")
	}

	return strings.TrimSpace(sb.String()), nil
}
