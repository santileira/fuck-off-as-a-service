package server

type Options struct {
	LogLevel                      string
	RateLimitCount                int
	RateLimitWindowInMilliseconds int
	TimeoutInMilliseconds         int
}
