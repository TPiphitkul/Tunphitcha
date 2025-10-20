package enforcement

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client

// InitRedis เริ่มเชื่อมต่อ Redis
func InitRedis(addr string) {
	rdb = redis.NewClient(&redis.Options{
		Addr: addr,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("Redis connection failed: %v", err))
	}
	fmt.Println("Connected to Redis:", addr)
}

// TokenBucketLimiter ใช้ Token-bucket ต่อ IP
func TokenBucketLimiter(ip string, rate int, capacity int) bool {
	key := fmt.Sprintf("bucket:%s", ip)
	// เติม token ทุก 1 วินาที
	now := time.Now().Unix()
	lastRefillKey := fmt.Sprintf("last:%s", ip)

	lastRefill, _ := rdb.Get(ctx, lastRefillKey).Int64()
	if lastRefill == 0 {
		lastRefill = now
	}

	elapsed := now - lastRefill
	if elapsed > 0 {
		// เติม token = elapsed * rate/60
		tokensToAdd := int(elapsed * int64(rate) / 60)
		current, _ := rdb.Get(ctx, key).Int()
		newTokens := current + tokensToAdd
		if newTokens > capacity {
			newTokens = capacity
		}
		rdb.Set(ctx, key, newTokens, 0)
		rdb.Set(ctx, lastRefillKey, now, 0)
	}

	// ใช้ token 1 อัน
	current, _ := rdb.Get(ctx, key).Int()
	if current <= 0 {
		return false
	}
	rdb.Decr(ctx, key)
	return true
}

// CountReqPerMin นับจำนวนคำขอต่อ IP (expire 60 วินาที)
func CountReqPerMin(ip string) int {
	key := fmt.Sprintf("reqpm:%s", ip)
	rdb.Incr(ctx, key)
	rdb.Expire(ctx, key, 60*time.Second)
	val, _ := rdb.Get(ctx, key).Int()
	return val
}

// Apply Middleware บังคับ rate limit
func Apply(decisionRate int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.Header.Get("X-Forwarded-For")
			if ip == "" {
				ip = r.RemoteAddr
			}
			ok := TokenBucketLimiter(ip, decisionRate, decisionRate)
			reqCount := CountReqPerMin(ip)
			w.Header().Set("X-ReqPerMin", strconv.Itoa(reqCount))
			w.Header().Set("X-Effective-RateLimit", strconv.Itoa(decisionRate))

			if !ok {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
