package converter

import (
	"context"
	"strings"
	"testing"
)

func TestDataConverter_CSV(t *testing.T) {
	csvData := `Name,Age
Alice,28
Bob,34`

	converter := &DataConverter{}
	opts := &Options{Extension: ".csv"}
	
	res, err := converter.Convert(context.Background(), strings.NewReader(csvData), opts)
	if err != nil {
		t.Fatalf("Failed to convert CSV: %v", err)
	}

	expected := "| Name | Age |\n|---|---|\n| Alice | 28 |\n| Bob | 34 |\n"
	if res != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, res)
	}
}
