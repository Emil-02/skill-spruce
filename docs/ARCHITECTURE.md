# Arsitektur

`skill-spruce` adalah satu binary Go tanpa dependensi eksternal. Ia hanya membaca file dan menulis laporan.

```
cmd/skill-spruce   → CLI (flag standar: check, rules, version)
internal/scan      → cari file SKILL.md / AGENTS.md / GEMINI.md / .cursor/rules/*.mdc
internal/rules     → tujuh aturan + interface Rule
internal/score     → agregasi bobot, grade, pembatasan skor kritis
internal/report    → formatter text / json / sarif / badge
internal/config    → konfigurasi opsional (rencana Fase 2)
```

## Kontrak aturan

```go
type Rule interface {
    ID() string
    Weight() int
    Severity() config.Severity
    Run(ctx context.Context, doc scan.Document) Result
}
```

Aturan baru tinggal mengimplementasikan interface ini dan didaftarkan di `rules.All()`. Tidak perlu mengubah pemindai maupun pelapor.

## Dokumen yang dipindai

`scan.Document` berisi `Path` (relatif), `Frontmatter`, dan `Body`. Frontmatter dipisahkan dari isi dengan delimiter `---` pada awal file.

## Perhitungan skor

- `pass` → dapat bobot penuh
- `warn` → dapat setengah bobot
- `fail` → dapat nol
- `skip` → tidak dihitung dalam total

`skor = round(earned / total × 100)`. Jika ada `fail` dengan severity `critical`, skor dibatasi maksimal 40.

## Keamanan

Tidak ada jaringan, tidak ada eksekusi perintah dari konten skill, tidak ada penulisan file. File yang mengandung pola berbahaya hanya dilaporkan, bukan dijalankan.
