package config

type Severity string

const (
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)

var rank = map[Severity]int{
	SeverityLow: 1, SeverityMedium: 2, SeverityHigh: 3, SeverityCritical: 4,
}

func ParseSeverity(s string) (Severity, bool) {
	v := Severity(s)
	_, ok := rank[v]
	return v, ok
}

func (s Severity) Rank() int { return rank[s] }

type Config struct {
	FailOn Severity
}

func Default() Config {
	return Config{FailOn: SeverityHigh}
}
