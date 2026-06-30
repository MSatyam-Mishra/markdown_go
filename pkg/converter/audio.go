package converter

import (
	"context"
	"io"
)

type AudioConverter struct{}

func (c *AudioConverter) Convert(ctx context.Context, r io.Reader, opts *Options) (string, error) {
	// Stub for audio transcription.
	// In a full implementation, you would use OpenAI Whisper API or whisper.cpp here.
	return "### Audio File\n\n*(Transcription stub - Requires Whisper API or whisper.cpp integration)*", nil
}
