package limiter

const (
	FixedWindow = iota
	SlidingWindow
	TokenBucket
	LeakyBucket
)
