package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/Emil-02/skill-spruce/internal/rules"
	"github.com/Emil-02/skill-spruce/internal/score"
)

func sample() Payload {
	return Payload{
		Root: "/tmp/demo",
		Report: score.Report{
			Score: 88, Grade: "A", Files: []string{"SKILL.md"},
			Results: []rules.Result{{
				ID: "description", Title: "Deskripsi spesifik", Status: "warn",
				Message: "deskripsi terlalu umum", File: "SKILL.md",
			}},
		},
	}
}

func TestRender(t *testing.T) {
	for _, format := range []string{"text", "json", "sarif", "badge"} {
		var buf bytes.Buffer
		if err := Render(&buf, sample(), format); err != nil {
			t.Fatalf("%s: %v", format, err)
		}
		if buf.Len() == 0 {
			t.Fatalf("%s kosong", format)
		}
		if format == "json" || format == "sarif" {
			var v any
			if err := json.Unmarshal(buf.Bytes(), &v); err != nil {
				t.Fatalf("%s bukan JSON: %v", format, err)
			}
		}
		if format == "badge" && !strings.Contains(buf.String(), "<svg") {
			t.Fatal("bukan SVG")
		}
	}
}

func TestUnknown(t *testing.T) {
	var buf bytes.Buffer
	if err := Render(&buf, sample(), "xml"); err == nil {
		t.Fatal("format asing harus error")
	}
}
