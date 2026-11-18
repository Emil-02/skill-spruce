package rules

import (
	"context"
	"regexp"
	"strings"

	"github.com/Emil-02/skill-spruce/internal/config"
	"github.com/Emil-02/skill-spruce/internal/scan"
)

type Result struct {
	ID       string          `json:"id"`
	Title    string          `json:"title"`
	Status   string          `json:"status"`
	Severity config.Severity `json:"severity"`
	Weight   int             `json:"weight"`
	Message  string          `json:"message,omitempty"`
	Hint     string          `json:"hint,omitempty"`
	File     string          `json:"file,omitempty"`
}

type Rule interface {
	ID() string
	Title() string
	Weight() int
	Severity() config.Severity
	Run(ctx context.Context, doc scan.Document) Result
}

type base struct {
	id, title string
	weight    int
	severity  config.Severity
}

func (b base) ID() string                { return b.id }
func (b base) Title() string             { return b.title }
func (b base) Weight() int               { return b.weight }
func (b base) Severity() config.Severity { return b.severity }

func skipNonSkill(r Rule, doc scan.Document) (Result, bool) {
	if doc.Kind == "" || doc.Kind == "skill" {
		return Result{}, false
	}
	return result(r, doc, "skip", "hanya berlaku untuk SKILL.md", ""), true
}

func result(r Rule, doc scan.Document, status, message, hint string) Result {
	return Result{
		ID: r.ID(), Title: r.Title(), Status: status, Severity: r.Severity(),
		Weight: r.Weight(), Message: message, Hint: hint, File: doc.Path,
	}
}

type Frontmatter struct{ base }

func (r Frontmatter) Run(_ context.Context, doc scan.Document) Result {
	if res, ok := skipNonSkill(r, doc); ok {
		return res
	}
	if strings.TrimSpace(doc.Frontmatter) == "" {
		return result(r, doc, "fail", "tidak ada frontmatter", "tambahkan blok --- dengan name dan description")
	}
	fm := strings.ToLower(doc.Frontmatter)
	missing := []string{}
	if !strings.Contains(fm, "name:") {
		missing = append(missing, "name")
	}
	if !strings.Contains(fm, "description:") {
		missing = append(missing, "description")
	}
	if len(missing) > 0 {
		return result(r, doc, "fail", "frontmatter kurang: "+strings.Join(missing, ", "), "lengkapi field tersebut")
	}
	return result(r, doc, "pass", "", "")
}

type Description struct{ base }

func (r Description) Run(_ context.Context, doc scan.Document) Result {
	if res, ok := skipNonSkill(r, doc); ok {
		return res
	}
	fm := strings.ToLower(doc.Frontmatter)
	desc := strings.ToLower(doc.Frontmatter)
	for _, key := range []string{"name:", "description:", "license:", "allowed-tools:", ">", "|"} {
		desc = strings.ReplaceAll(desc, key, " ")
	}
	desc = strings.Join(strings.Fields(desc), " ")
	if len(desc) < 40 {
		return result(r, doc, "warn", "deskripsi terlalu pendek", "jelaskan apa dan kapan skill dipakai")
	}
	for _, bad := range []string{"replace", "todo", "deskripsi di sini", "replace with"} {
		if strings.Contains(fm, bad) {
			return result(r, doc, "fail", "deskripsi masih placeholder", "tulis deskripsi yang spesifik")
		}
	}
	triggers := []string{"when", "kapan", "use", "pakai", "sebelum", "saat", "jika", "gunakan"}
	for _, t := range triggers {
		if strings.Contains(fm, t) {
			return result(r, doc, "pass", "", "")
		}
	}
	return result(r, doc, "warn", "deskripsi tidak menyebut pemicu", "tambahkan kalimat 'Gunakan skill ini saat ...'")
}

var (
	curlBash  = regexp.MustCompile(`(?i)curl[^\n]*\|\s*(sudo\s*)?bash`)
	rmRoot    = regexp.MustCompile(`(?i)rm\s+-rf\s+/(\s|$)`)
	printKeys = regexp.MustCompile(`(?i)(echo|print|cat)\s+[^\n]*(api[_-]?key|secret|token|password)`)
	evalCall  = regexp.MustCompile(`(?i)eval\s*\(`)
	autoPush  = regexp.MustCompile(`(?i)auto[ -]?push`)
	negations = []string{"jangan", "dilarang", "hindari", "never", "must not", "do not", "tidak boleh"}
)

type Dangerous struct{ base }

func (r Dangerous) Run(_ context.Context, doc scan.Document) Result {
	patterns := []*regexp.Regexp{curlBash, rmRoot, printKeys, evalCall, autoPush}
	labels := []string{"curl | bash", "rm -rf /", "print secret", "eval(", "auto-push"}
	for _, line := range strings.Split(doc.Body, "\n") {
		lower := strings.ToLower(line)
		negated := false
		for _, n := range negations {
			if strings.Contains(lower, n) {
				negated = true
				break
			}
		}
		if negated {
			continue
		}
		for i, p := range patterns {
			if p.MatchString(line) {
				return result(r, doc, "fail", "instruksi berbahaya: "+labels[i], "hapus atau beri batas eksplisit")
			}
		}
	}
	return result(r, doc, "pass", "", "")
}

type HardRules struct{ base }

func (r HardRules) Run(_ context.Context, doc scan.Document) Result {
	count := 0
	for _, line := range strings.Split(strings.ToLower(doc.Body), "\n") {
		for _, kw := range []string{"jangan", "dilarang", "tidak boleh", "never", "must not"} {
			if strings.Contains(line, kw) {
				count++
				break
			}
		}
	}
	switch {
	case count >= 2:
		return result(r, doc, "pass", "", "")
	case count == 1:
		return result(r, doc, "warn", "baru ada satu aturan keras", "tambahkan batas soal secret dan auto-push")
	default:
		return result(r, doc, "fail", "tidak ada aturan keras", "tulis setidaknya dua aturan: jangan bocorkan secret, jangan auto-push")
	}
}

type Length struct{ base }

func (r Length) Run(_ context.Context, doc scan.Document) Result {
	n := len(doc.Body)
	switch {
	case n > 8000 && !strings.Contains(doc.Body, "references/"):
		return result(r, doc, "fail", "instruksi terlalu panjang", "pindahkan detail ke references/")
	case n > 8000:
		return result(r, doc, "warn", "instruksi panjang", "pertimbangkan memendekkan")
	default:
		return result(r, doc, "pass", "", "")
	}
}

type AgentCoverage struct{ base }

func (r AgentCoverage) Run(_ context.Context, doc scan.Document) Result {
	all := strings.ToLower(doc.Frontmatter + "\n" + doc.Body)
	names := []string{"opencode", "claude", "codex", "gemini", "cursor", "agents.md"}
	found := 0
	for _, n := range names {
		if strings.Contains(all, n) {
			found++
		}
	}
	switch {
	case found >= 2:
		return result(r, doc, "pass", "", "")
	case found == 1:
		return result(r, doc, "warn", "baru ditujukan untuk satu agent", "sebutkan agent lain yang didukung")
	default:
		return result(r, doc, "fail", "tidak menyebut agent mana pun", "tambahkan daftar agent yang didukung")
	}
}

type Progressive struct{ base }

func (r Progressive) Run(_ context.Context, doc scan.Document) Result {
	if res, ok := skipNonSkill(r, doc); ok {
		return res
	}
	if strings.Contains(doc.Body, "references/") {
		return result(r, doc, "pass", "", "")
	}
	if len(doc.Body) > 6000 {
		return result(r, doc, "warn", "detail tidak dipisah ke references/", "pindahkan aturan panjang ke references/")
	}
	return result(r, doc, "pass", "", "")
}

func All() []Rule {
	return []Rule{
		Frontmatter{base{"frontmatter", "Frontmatter valid", 20, config.SeverityHigh}},
		Description{base{"description", "Deskripsi spesifik", 20, config.SeverityHigh}},
		Dangerous{base{"dangerous", "Tanpa instruksi berbahaya", 20, config.SeverityCritical}},
		HardRules{base{"hard-rules", "Aturan keras untuk agent", 15, config.SeverityHigh}},
		Length{base{"length", "Panjang wajar", 10, config.SeverityMedium}},
		AgentCoverage{base{"agent-coverage", "Cakupan multi-agent", 10, config.SeverityLow}},
		Progressive{base{"progressive", "Progressive disclosure", 5, config.SeverityLow}},
	}
}
