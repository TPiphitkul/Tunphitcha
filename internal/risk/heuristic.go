package risk

import (
	"github.com/example/go-adaptive-gw/internal/profiler"
)

// Score returns an int risk score (0-100) based on simple heuristics.
func Score(m profiler.Meta) int {
	score := 0
	// Sensitive endpoint heuristic
	if m.Path == "/api/user/login" && m.Method == "POST" {
		score += 20
	}
	// Burst heuristic
	if m.ReqPerMin > 100 {
		score += 40
	}
	// Unknown UA heuristic (very simple)
	if m.UserAgent == "" {
		score += 10
	}
	return score
}

func Level(score int) string {
	switch {
	case score >= 70:
		return "high"
	case score >= 40:
		return "medium"
	default:
		return "low"
	}
}
