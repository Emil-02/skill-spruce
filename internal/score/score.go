package score

import (
	"math"
	"sort"

	"github.com/Emil-02/skill-spruce/internal/config"
	"github.com/Emil-02/skill-spruce/internal/rules"
	"github.com/Emil-02/skill-spruce/internal/scan"
)

type Report struct {
	Score   int            `json:"score"`
	Grade   string         `json:"grade"`
	Files   []string       `json:"files"`
	Results []rules.Result `json:"results"`
}

var statusWeight = map[string]int{"skip": 0, "pass": 1, "warn": 2, "fail": 3}

func Evaluate(docs []scan.Document, cfg config.Config) Report {
	combined := map[string]rules.Result{}
	for _, doc := range docs {
		for _, rule := range rules.All() {
			res := rule.Run(nil, doc)
			old, ok := combined[res.ID]
			if !ok || statusWeight[res.Status] > statusWeight[old.Status] {
				combined[res.ID] = res
			}
		}
	}

	ids := make([]string, 0, len(combined))
	for id := range combined {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	rep := Report{Results: []rules.Result{}, Files: []string{}}
	order := map[string]int{}
	for i, r := range rules.All() {
		order[r.ID()] = i
	}
	sort.Slice(ids, func(i, j int) bool { return order[ids[i]] < order[ids[j]] })

	var earned, total float64
	critical := false
	for _, id := range ids {
		res := combined[id]
		rep.Results = append(rep.Results, res)
		total += float64(res.Weight)
		switch res.Status {
		case "pass":
			earned += float64(res.Weight)
		case "warn":
			earned += float64(res.Weight) / 2
		case "fail":
			if res.Severity == config.SeverityCritical {
				critical = true
			}
		}
	}
	for _, doc := range docs {
		rep.Files = append(rep.Files, doc.Path)
	}

	switch {
	case total == 0:
		rep.Score = 0
	default:
		rep.Score = int(math.Round(earned / total * 100))
	}
	if critical && rep.Score > 40 {
		rep.Score = 40
	}
	rep.Grade = Grade(rep.Score)
	return rep
}

func Grade(score int) string {
	switch {
	case score >= 95:
		return "A+"
	case score >= 85:
		return "A"
	case score >= 70:
		return "B"
	case score >= 55:
		return "C"
	case score >= 40:
		return "D"
	}
	return "F"
}

func Breached(results []rules.Result, failOn config.Severity) bool {
	for _, r := range results {
		if r.Status == "fail" && r.Severity.Rank() >= failOn.Rank() {
			return true
		}
	}
	return false
}
