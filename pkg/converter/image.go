package converter

import (
	"context"
	"fmt"
	"io"
	"strings"

	exif "github.com/dsoprea/go-exif/v3"
)

type ImageConverter struct{}

func (c *ImageConverter) Convert(ctx context.Context, r io.Reader, opts *Options) (string, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("## Image: %s\n\n", opts.FileName))

	// EXIF
	rawExif, err := exif.SearchAndExtractExif(b)
	if err == nil {
		sb.WriteString("### EXIF Metadata\n\n```json\n")
		sb.WriteString(fmt.Sprintf("Raw EXIF bytes size: %d\n", len(rawExif)))
		sb.WriteString("```\n\n")
	}

	// OCR
    sb.WriteString("### Text Content (OCR)\n\n")
    sb.WriteString("*(OCR transcription stub - To enable full OCR, integrate github.com/otiai10/gosseract/v2 or a cloud vision API)*\n")

	return sb.String(), nil
}
