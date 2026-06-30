package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/user/markdown_go"
)

func main() {
	var input string
	var output string
	var urlStr string

	flag.StringVar(&input, "i", "", "Input file path")
	flag.StringVar(&output, "o", "", "Output file path (optional, defaults to stdout)")
	flag.StringVar(&urlStr, "u", "", "Input URL (e.g., YouTube URL)")
	flag.Parse()

	if input == "" && urlStr == "" {
		flag.Usage()
		os.Exit(1)
	}

	m := markdown_go.NewMarkItDown()
	ctx := context.Background()
	var res string
	var err error

	if urlStr != "" {
		res, err = m.ConvertURL(ctx, urlStr)
	} else if input != "" {
		res, err = m.ConvertFile(ctx, input)
	}

	if err != nil {
		log.Fatalf("Conversion failed: %v\n", err)
	}

	if output != "" {
		err = os.WriteFile(output, []byte(res), 0644)
		if err != nil {
			log.Fatalf("Failed to write output file: %v\n", err)
		}
		fmt.Printf("Successfully converted to %s\n", output)
	} else {
		fmt.Println(res)
	}
}
