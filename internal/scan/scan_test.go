package scan

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Emil-02/skill-spruce/internal/config"
)

func TestFindAndSplit(t *testing.T) {
	root := t.TempDir()
	files := map[string]string{
		"SKILL.md":                "---\nname: demo\ndescription: Gunakan saat commit.\n---\n# Demo\n",
		"AGENTS.md":               "# Agent\n",
		".cursor/rules/demo.mdc":  "---\ndescription: x\n---\nrule\n",
		"examples/bad/SKILL.md":   "# jangan\n",
		"node_modules/x/SKILL.md": "# x\n",
	}
	for p, c := range files {
		full := filepath.Join(root, p)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(c), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	docs, err := Find(root, config.Default())
	if err != nil {
		t.Fatal(err)
	}
	if len(docs) != 3 {
		t.Fatalf("dapat %d dokumen: %+v", len(docs), docs)
	}
	for _, d := range docs {
		if d.Path == "examples/bad/SKILL.md" || d.Path == "node_modules/x/SKILL.md" {
			t.Fatal("direktori yang dilewati tetap dipindai:", d.Path)
		}
	}
	var skill Document
	for _, d := range docs {
		if d.Path == "SKILL.md" {
			skill = d
		}
	}
	if skill.Frontmatter == "" || skill.Body == "" {
		t.Fatalf("frontmatter/body salah: %+v", skill)
	}
}
