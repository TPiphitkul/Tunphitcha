package profiler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client

// Meta holds basic request metadata used for heuristic or ML scoring.
type Meta struct {
	IP         string `json:"ip"`
	UserAgent  string `json:"user_agent"`
	Path       string `json:"path"`
	Method     string `json:"method"`
	ReqPerMin  int    `json:"req_per_min"`
	IsAttack   bool   `json:"is_attack,omitempty"` // label สำหรับเก็บ dataset
	RiskScore  int    `json:"risk_score,omitempty"`
	RecordedAt string `json:"recorded_at,omitempty"`
}

// InitRedis initializes Redis connection for profiling and counting requests.
func InitRedis(addr string) {
	rdb = redis.NewClient(&redis.Options{
		Addr: addr,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic("Redis not reachable: " + err.Error())
	}
}

// Extract produces Meta from a request and counts requests per minute (via Redis).
func Extract(r *http.Request) Meta {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	ua := strings.TrimSpace(r.UserAgent())

	// Count request per minute
	key := "reqpm:" + ip
	rdb.Incr(ctx, key)
	rdb.Expire(ctx, key, 60*time.Second)
	reqPerMin, _ := rdb.Get(ctx, key).Int()

	return Meta{
		IP:         ip,
		UserAgent:  ua,
		Path:       r.URL.Path,
		Method:     r.Method,
		ReqPerMin:  reqPerMin,
		RecordedAt: time.Now().Format(time.RFC3339),
	}
}

// Export writes request metadata to a JSON log for dataset building (offline training).
func Export(m Meta) {
	data, _ := json.Marshal(m)
	rdb.LPush(ctx, "risk_log", data)
	// หรือถ้าอยากเก็บลงไฟล์เพิ่ม:
	// f, _ := os.OpenFile("risk_log.jsonl", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// defer f.Close()
	// f.Write(append(data, '\n'))
}
