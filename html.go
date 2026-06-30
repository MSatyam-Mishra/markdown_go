package markdown_go

import (
	"context"
	"io"
	"strings"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown"
)

type HTMLConverter struct{}

func (c *HTMLConverter) Convert(ctx context.Context, r io.Reader, opts *Options) (string, error) {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	converter := htmltomarkdown.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(string(bytes))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(markdown), nil
}
