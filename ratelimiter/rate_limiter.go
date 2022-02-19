package ratelimiter

type RateLimiter interface {
	AllowRequest(userID string) bool
}
