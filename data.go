package markdown_go

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type DataConverter struct{}

func (c *DataConverter) Convert(ctx context.Context, r io.Reader, opts *Options) (string, error) {
	if opts.Extension == ".csv" {
		return c.convertCSV(r)
	} else if opts.Extension == ".json" {
		return c.convertJSON(r)
	} else if opts.Extension == ".xml" {
		// return as code block
		b, err := io.ReadAll(r)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("```xml\n%s\n```", string(b)), nil
	}
	return "", fmt.Errorf("unsupported data extension: %s", opts.Extension)
}

func (c *DataConverter) convertCSV(r io.Reader) (string, error) {
	reader := csv.NewReader(r)
	records, err := reader.ReadAll()
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	for i, row := range records {
		sb.WriteString("|")
		for _, col := range row {
			sb.WriteString(fmt.Sprintf(" %s |", strings.ReplaceAll(col, "\n", " ")))
		}
		sb.WriteString("\n")

		if i == 0 {
			sb.WriteString("|")
			for range row {
				sb.WriteString("---|")
			}
			sb.WriteString("\n")
		}
	}
	return sb.String(), nil
}

func (c *DataConverter) convertJSON(r io.Reader) (string, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	// Try pretty printing json
	var obj interface{}
	if err := json.Unmarshal(b, &obj); err == nil {
		pretty, err := json.MarshalIndent(obj, "", "  ")
		if err == nil {
			return fmt.Sprintf("```json\n%s\n```", string(pretty)), nil
		}
	}
	return fmt.Sprintf("```json\n%s\n```", string(b)), nil
}
