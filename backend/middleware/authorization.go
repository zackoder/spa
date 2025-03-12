package middleware

import (
	"net/http"
	"time"

	"reat-time-forum/utils"
)

type Limit struct {
	LastTime int64
	Counter  int
}

type RateLimit struct {
	User map[string]Limit
}

func (r *RateLimit) Allow(ip string) bool {
	if _, ok := r.User[ip]; !ok {
		currentTime := time.Now()
		r.User[ip] = Limit{LastTime: currentTime.Unix(), Counter: 1}
	} else {
		limit := r.User[ip]
		currentTime := time.Now()
		if limit.LastTime == currentTime.Unix() {
			limit.Counter++
			if limit.Counter > 20 {
				return false
			}
		} else {
			limit.LastTime = currentTime.Unix()
			limit.Counter = 1
		}
	}
	return true
}

var rateLimit RateLimit

func Authorization(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie := utils.CheckCookie(r)
		if cookie == nil {
			w.WriteHeader(401)
			return
		}
		allowed := rateLimit.Allow(r.RemoteAddr)
		if !allowed {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}
