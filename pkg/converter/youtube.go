package converter

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/chand1012/yt_transcript"
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

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s\n\n**Author**: %s\n**Duration**: %s\n\n### Description\n%s\n\n",
		video.Title, video.Author, video.Duration.String(), video.Description))

	// Fetch Captions
	transcripts, _, err := yt_transcript.FetchTranscript(video.ID, "en", "US")
	if err == nil && len(transcripts) > 0 {
		sb.WriteString("### Transcript\n\n")
		for _, t := range transcripts {
			sb.WriteString(t.Text)
			sb.WriteString(" ") // use spaces for flow
		}
		sb.WriteString("\n")
	} else {
		sb.WriteString("*(No English transcript available)*\n")
	}

	return sb.String(), nil
}
