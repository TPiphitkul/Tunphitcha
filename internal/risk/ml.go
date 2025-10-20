package risk

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
	"github.com/example/go-adaptive-gw/internal/profiler"
)

type MLResponse struct {
	RiskScore int `json:"risk_score"`
}

func MLScore(m profiler.Meta) int {
	payload, _ := json.Marshal(map[string]interface{}{
		"req_per_min": m.ReqPerMin,
		"path":        m.Path,
	})
	client := &http.Client{Timeout: 500 * time.Millisecond}
	resp, err := client.Post("http://ml-risk:5000/predict", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return 50 // fallback
	}
	defer resp.Body.Close()
	var result MLResponse
	json.NewDecoder(resp.Body).Decode(&result)
	return result.RiskScore
}
