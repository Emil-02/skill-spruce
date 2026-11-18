package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Emil-02/skill-spruce/internal/config"
	"github.com/Emil-02/skill-spruce/internal/report"
	"github.com/Emil-02/skill-spruce/internal/rules"
	"github.com/Emil-02/skill-spruce/internal/scan"
	"github.com/Emil-02/skill-spruce/internal/score"
)

const version = "0.1.0"

type options struct {
	path   string
	format string
	failOn string
	output string
	ci     bool
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "check":
		opts, err := parse(os.Args[2:])
		if err != nil {
			fmt.Fprintln(os.Stderr, "skill-spruce:", err)
			os.Exit(2)
		}
		os.Exit(runCheck(opts))
	case "rules":
		for _, r := range rules.All() {
			fmt.Printf("%-16s %-3d %-8s %s\n", r.ID(), r.Weight(), r.Severity(), r.Title())
		}
	case "version":
		fmt.Println("skill-spruce", version)
	case "help", "-h", "--help":
		usage()
	default:
		fmt.Fprintln(os.Stderr, "skill-spruce: perintah tidak dikenal:", os.Args[1])
		usage()
		os.Exit(2)
	}
}

func parse(args []string) (options, error) {
	opts := options{path: ".", format: "text"}
	for i := 0; i < len(args); i++ {
		a := args[i]
		next := func() (string, error) {
			if i+1 >= len(args) {
				return "", fmt.Errorf("nilai untuk %s tidak ada", a)
			}
			i++
			return args[i], nil
		}
		switch a {
		case "--format":
			v, err := next()
			if err != nil {
				return opts, err
			}
			opts.format = v
		case "--fail-on":
			v, err := next()
			if err != nil {
				return opts, err
			}
			opts.failOn = v
		case "-o", "--output":
			v, err := next()
			if err != nil {
				return opts, err
			}
			opts.output = v
		case "--ci":
			opts.ci = true
		default:
			if opts.path != "." {
				return opts, fmt.Errorf("argumen tidak dikenal: %s", a)
			}
			opts.path = a
		}
	}
	return opts, nil
}

func runCheck(opts options) int {
	root, err := filepath.Abs(opts.path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "skill-spruce:", err)
		return 2
	}
	info, err := os.Stat(root)
	if err != nil || !info.IsDir() {
		fmt.Fprintln(os.Stderr, "skill-spruce: bukan direktori:", opts.path)
		return 2
	}

	cfg := config.Default()
	if opts.failOn != "" {
		sev, ok := config.ParseSeverity(opts.failOn)
		if !ok {
			fmt.Fprintln(os.Stderr, "skill-spruce: fail_on tidak valid:", opts.failOn)
			return 2
		}
		cfg.FailOn = sev
	}

	docs, err := scan.Find(root, cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "skill-spruce:", err)
		return 2
	}
	rep := score.Evaluate(docs, cfg)

	var out io.Writer = os.Stdout
	if opts.output != "" {
		f, err := os.Create(opts.output)
		if err != nil {
			fmt.Fprintln(os.Stderr, "skill-spruce:", err)
			return 2
		}
		defer f.Close()
		out = f
	}
	if err := report.Render(out, report.Payload{Root: root, Report: rep}, opts.format); err != nil {
		fmt.Fprintln(os.Stderr, "skill-spruce:", err)
		return 2
	}
	if opts.ci && score.Breached(rep.Results, cfg.FailOn) {
		return 1
	}
	return 0
}

func usage() {
	fmt.Print(`skill-spruce — audit kualitas skill untuk AI coding agent

Pemakaian:
  skill-spruce check [path] [opsi]
  skill-spruce rules
  skill-spruce version

Opsi:
  --format text|json|sarif|badge   format keluaran (default: text)
  -o, --output FILE                tulis ke file
  --fail-on low|medium|high|critical  ambang kegagalan
  --ci                             keluar kode 1 bila ada temuan fatal

Contoh:
  skill-spruce check .
  skill-spruce check --format badge -o badge.svg .
  skill-spruce check --ci --fail-on high .
`)
}
