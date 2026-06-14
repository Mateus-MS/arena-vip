package contato

import (
	"sync"
	"time"
)

const (
	rateWindow  = 10 * 60 // 10-minute sliding window (seconds)
	rateMaxHits = 3       // max submissions per window per IP
)

var (
	rateMu  sync.Mutex
	rateMap = make(map[string][]int64)
)

func isRateLimited(ip string) bool {
	rateMu.Lock()
	defer rateMu.Unlock()

	now := time.Now().Unix()
	hits := rateMap[ip]

	var recent []int64
	for _, t := range hits {
		if now-t < rateWindow {
			recent = append(recent, t)
		}
	}

	if len(recent) >= rateMaxHits {
		rateMap[ip] = recent
		return true
	}

	rateMap[ip] = append(recent, now)
	return false
}
