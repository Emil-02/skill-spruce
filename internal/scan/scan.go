package scan

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/Emil-02/skill-spruce/internal/config"
)

type Document struct {
	Path        string `json:"path"`
	Kind        string `json:"kind"`
	Frontmatter string `json:"frontmatter"`
	Body        string `json:"body"`
}

var skipDirs = map[string]bool{
	".git": true, "node_modules": true, "dist": true, "build": true,
	"vendor": true, "coverage": true, "testdata": true, "examples": true,
}

func Find(root string, cfg config.Config) ([]Document, error) {
	var docs []Document
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			hidden := strings.HasPrefix(d.Name(), ".") && d.Name() != "."
			allowedHidden := d.Name() == ".cursor" || d.Name() == ".opencode" || d.Name() == ".claude" || d.Name() == ".gemini"
			if path != root && (skipDirs[d.Name()] || hidden && !allowedHidden) {
				return filepath.SkipDir
			}
			return nil
		}
		if !isSkillFile(path) {
			return nil
		}
		info, err := d.Info()
		if err != nil || info.Size() > 200<<10 {
			return nil
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		text := string(raw)
		if strings.Contains(text, "BEGIN:nextjs-agent-rules") {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			rel = path
		}
		kind := "agents"
		if strings.EqualFold(d.Name(), "skill.md") {
			kind = "skill"
		}
		fm, body := splitFrontmatter(text)
		docs = append(docs, Document{Path: filepath.ToSlash(rel), Kind: kind, Frontmatter: fm, Body: body})
		return nil
	})
	return docs, err
}

func isSkillFile(path string) bool {
	base := strings.ToLower(filepath.Base(path))
	normalized := filepath.ToSlash(path)
	switch base {
	case "skill.md", "agents.md", "gemini.md", "claude.md":
		return true
	}
	return strings.Contains(normalized, "/.cursor/rules/") && strings.HasSuffix(base, ".mdc")
}

func splitFrontmatter(raw string) (string, string) {
	if !strings.HasPrefix(raw, "---") {
		return "", raw
	}
	lines := strings.Split(raw, "\n")
	if len(lines) < 2 || strings.TrimSpace(lines[0]) != "---" {
		return "", raw
	}
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			return strings.Join(lines[1:i], "\n"), strings.Join(lines[i+1:], "\n")
		}
	}
	return "", raw
}
