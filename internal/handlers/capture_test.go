package handlers

import (
	"testing"

	"pagemail/internal/models"
)

func TestFormatsToInt(t *testing.T) {
	tests := []struct {
		name    string
		formats []string
		want    int
	}{
		{"empty", []string{}, 0},
		{"pdf only", []string{"pdf"}, models.FormatPDF},
		{"html only", []string{"html"}, models.FormatHTML},
		{"screenshot only", []string{"screenshot"}, models.FormatPNG},
		{"all formats", []string{"pdf", "html", "screenshot"}, models.FormatPDF | models.FormatHTML | models.FormatPNG},
		{"pdf and html", []string{"pdf", "html"}, models.FormatPDF | models.FormatHTML},
		{"unknown format ignored", []string{"pdf", "unknown"}, models.FormatPDF},
		{"duplicates", []string{"pdf", "pdf"}, models.FormatPDF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatsToInt(tt.formats)
			if got != tt.want {
				t.Errorf("formatsToInt(%v) = %d, want %d", tt.formats, got, tt.want)
			}
		})
	}
}

func TestIntToFormats(t *testing.T) {
	tests := []struct {
		name  string
		flags int
		want  []string
	}{
		{"zero", 0, nil},
		{"pdf only", models.FormatPDF, []string{"pdf"}},
		{"html only", models.FormatHTML, []string{"html"}},
		{"screenshot only", models.FormatPNG, []string{"screenshot"}},
		{"all formats", models.FormatPDF | models.FormatHTML | models.FormatPNG, []string{"pdf", "html", "screenshot"}},
		{"pdf and screenshot", models.FormatPDF | models.FormatPNG, []string{"pdf", "screenshot"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := intToFormats(tt.flags)
			if len(got) != len(tt.want) {
				t.Errorf("intToFormats(%d) length = %d, want %d", tt.flags, len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("intToFormats(%d)[%d] = %q, want %q", tt.flags, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestFormatRoundTrip(t *testing.T) {
	testCases := [][]string{
		{"pdf"},
		{"html"},
		{"screenshot"},
		{"pdf", "html"},
		{"pdf", "screenshot"},
		{"html", "screenshot"},
		{"pdf", "html", "screenshot"},
	}

	for _, formats := range testCases {
		t.Run("roundtrip", func(t *testing.T) {
			flags := formatsToInt(formats)
			result := intToFormats(flags)

			if len(result) != len(formats) {
				t.Errorf("Round trip failed: %v -> %d -> %v", formats, flags, result)
			}
		})
	}
}
