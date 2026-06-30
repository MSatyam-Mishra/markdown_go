package converter

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"strings"
)

type EpubConverter struct {
	HTMLConverter *HTMLConverter
}

func (c *EpubConverter) Convert(ctx context.Context, r io.Reader, opts *Options) (string, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	bytesReader := bytes.NewReader(b)

	zr, err := zip.NewReader(bytesReader, bytesReader.Size())
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	for _, f := range zr.File {
		lowerName := strings.ToLower(f.Name)
		if strings.HasSuffix(lowerName, ".html") || strings.HasSuffix(lowerName, ".htm") || strings.HasSuffix(lowerName, ".xhtml") {
			rc, err := f.Open()
			if err != nil {
				continue
			}

			res, err := c.HTMLConverter.Convert(ctx, rc, &Options{Extension: ".html"})
			rc.Close()
			if err == nil {
				sb.WriteString(res)
				sb.WriteString("\n\n---\n\n")
			}
		}
	}
	return sb.String(), nil
}
