package profiler

import (
	"net/http"
	"strings"
)

// Meta holds basic request metadata used for heuristic scoring.
type Meta struct {
	IP        string
	UserAgent string
	Path      string
	Method    string
	ReqPerMin int
}

// Extract produces Meta from a request.
func Extract(r *http.Request) Meta {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	ua := r.UserAgent()
	ua = strings.TrimSpace(ua)
	return Meta{
		IP:        ip,
		UserAgent: ua,
		Path:      r.URL.Path,
		Method:    r.Method,
	}
}
