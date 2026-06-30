package markdown_go

import (
	"context"
	"fmt"
	"io"

	"github.com/kkdai/youtube/v2"
)

type YoutubeConverter struct{}

func (c *YoutubeConverter) Convert(ctx context.Context, r io.Reader, opts *Options) (string, error) {
	if opts.URL == "" {
		return "", fmt.Errorf("no URL provided for Youtube converter")
	}

	client := youtube.Client{}
	video, err := client.GetVideo(opts.URL)
	if err != nil {
		return "", err
	}

	markdown := fmt.Sprintf("# %s\n\n**Author**: %s\n**Duration**: %s\n\n%s\n",
		video.Title, video.Author, video.Duration.String(), video.Description)

	return markdown, nil
}
