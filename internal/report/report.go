package report

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/Emil-02/skill-spruce/internal/score"
)

type Payload struct {
	Root   string       `json:"root"`
	Report score.Report `json:"report"`
}

func Render(w io.Writer, p Payload, format string) error {
	switch format {
	case "text":
		return renderText(w, p)
	case "json":
		return json.NewEncoder(w).Encode(p)
	case "sarif":
		return renderSARIF(w, p)
	case "badge":
		return renderBadge(w, p)
	default:
		return fmt.Errorf("format tidak dikenal: %s", format)
	}
}

func renderText(w io.Writer, p Payload) error {
	var b strings.Builder
	fmt.Fprintf(&b, "skill-spruce — %s\n", p.Root)
	b.WriteString(strings.Repeat("─", 56) + "\n")
	for _, r := range p.Report.Results {
		mark := "✓"
		if r.Status == "warn" {
			mark = "!"
		}
		if r.Status == "fail" {
			mark = "✗"
		}
		fmt.Fprintf(&b, "%s %-16s %s\n", mark, r.ID, strings.ToUpper(string(r.Severity)))
		if r.Message != "" {
			fmt.Fprintf(&b, "    %s: %s\n", r.File, r.Message)
		}
		if r.Hint != "" {
			fmt.Fprintf(&b, "    ↳ %s\n", r.Hint)
		}
	}
	b.WriteString(strings.Repeat("─", 56) + "\n")
	fmt.Fprintf(&b, "Skor: %d/100 (%s) — %d file skill\n", p.Report.Score, p.Report.Grade, len(p.Report.Files))
	_, err := io.WriteString(w, b.String())
	return err
}

type sarifResult struct {
	RuleID  string `json:"ruleId"`
	Level   string `json:"level"`
	Message struct {
		Text string `json:"text"`
	} `json:"message"`
	Locations []struct {
		PhysicalLocation struct {
			ArtifactLocation struct {
				URI string `json:"uri"`
			} `json:"artifactLocation"`
			Region struct {
				StartLine int `json:"startLine"`
			} `json:"region"`
		} `json:"physicalLocation"`
	} `json:"locations"`
}

type sarifRun struct {
	Tool struct {
		Driver struct {
			Name string `json:"name"`
		} `json:"driver"`
	} `json:"tool"`
	Results []sarifResult `json:"results"`
}

func renderSARIF(w io.Writer, p Payload) error {
	run := sarifRun{Results: []sarifResult{}}
	run.Tool.Driver.Name = "skill-spruce"
	for _, r := range p.Report.Results {
		if r.Status != "fail" && r.Status != "warn" {
			continue
		}
		item := sarifResult{RuleID: r.ID, Level: "warning"}
		if r.Status == "fail" {
			item.Level = "error"
		}
		item.Message.Text = r.Message
		loc := struct {
			PhysicalLocation struct {
				ArtifactLocation struct {
					URI string `json:"uri"`
				} `json:"artifactLocation"`
				Region struct {
					StartLine int `json:"startLine"`
				} `json:"region"`
			} `json:"physicalLocation"`
		}{}
		loc.PhysicalLocation.ArtifactLocation.URI = r.File
		loc.PhysicalLocation.Region.StartLine = 1
		item.Locations = append(item.Locations, loc)
		run.Results = append(run.Results, item)
	}
	doc := struct {
		Version string     `json:"version"`
		Schema  string     `json:"$schema"`
		Runs    []sarifRun `json:"runs"`
	}{
		Version: "2.1.0",
		Schema:  "https://json.schemastore.org/sarif-2.1.0.json",
		Runs:    []sarifRun{run},
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(doc)
}

func renderBadge(w io.Writer, p Payload) error {
	color := "#4c1"
	switch {
	case p.Report.Score < 50:
		color = "#e05d44"
	case p.Report.Score < 80:
		color = "#dfb317"
	}
	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="126" height="20" role="img" aria-label="skill health: %d/100">
  <linearGradient id="s" x2="0" y2="100%%"><stop offset="0" stop-color="#bbb" stop-opacity=".1"/><stop offset="1" stop-opacity=".1"/></linearGradient>
  <clipPath id="r"><rect width="126" height="20" rx="3" fill="#fff"/></clipPath>
  <g clip-path="url(#r)"><rect width="86" height="20" fill="#555"/><rect x="86" width="40" height="20" fill="%s"/><rect width="126" height="20" fill="url(#s)"/></g>
  <g fill="#fff" text-anchor="middle" font-family="Verdana,Geneva,DejaVu Sans,sans-serif" font-size="11">
    <text x="43" y="15">skill health</text><text x="106" y="15">%d/100</text>
  </g>
</svg>
`, p.Report.Score, color, p.Report.Score)
	_, err := io.WriteString(w, svg)
	return err
}
