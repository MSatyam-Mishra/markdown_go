package converter

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"sync"
)

type ZipConverter struct {
	API ConverterAPI
}

func (c *ZipConverter) Convert(ctx context.Context, r io.Reader, opts *Options) (string, error) {
	// ZIP requires an io.ReaderAt, so we must read the whole stream into memory first.
	b, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	bytesReader := bytes.NewReader(b)

	zr, err := zip.NewReader(bytesReader, bytesReader.Size())
	if err != nil {
		return "", err
	}

	var wg sync.WaitGroup
	results := make(chan string, len(zr.File))

	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			continue
		}

		wg.Add(1)
		go func(file *zip.File) {
			defer wg.Done()
			rc, err := file.Open()
			if err != nil {
				return
			}
			defer rc.Close()

			ext := strings.ToLower(filepath.Ext(file.Name))
			fileOpts := &Options{
				Extension: ext,
				FileName:  file.Name,
			}

			res, err := c.API.Convert(ctx, rc, fileOpts)
			if err != nil {
				// Ignore unsupported files in zip
				return
			}

			header := fmt.Sprintf("\n\n### File: %s\n\n", file.Name)
			results <- header + res
		}(f)
	}

	wg.Wait()
	close(results)

	var sb strings.Builder
	for res := range results {
		sb.WriteString(res)
	}

	return sb.String(), nil
}
