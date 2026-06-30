package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/MSatyam-Mishra/markdown_go"
	"github.com/MSatyam-Mishra/markdown_go/pkg/converter"
)

//go:embed index.html
var indexHTML []byte

func main() {
	m := markdown_go.NewMarkItDown()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(indexHTML)
	})

	http.HandleFunc("/api/convert", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseMultipartForm(32 << 20) // 32MB max
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Handle file upload
		file, handler, err := r.FormFile("file")
		if err == nil {
			defer file.Close()
			ext := strings.ToLower(filepath.Ext(handler.Filename))
			opts := &converter.Options{
				Extension: ext,
				FileName:  handler.Filename,
			}

			md, err := m.Convert(context.Background(), file, opts)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"markdown": md})
			return
		}

		// Handle URL upload (like Youtube)
		urlStr := r.FormValue("url")
		if urlStr != "" {
			md, err := m.ConvertURL(context.Background(), urlStr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"markdown": md})
			return
		}

		http.Error(w, "No file or URL provided", http.StatusBadRequest)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Frontend is running at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
