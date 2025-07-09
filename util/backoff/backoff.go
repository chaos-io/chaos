package backoff

import (
	"math"
	"math/rand"
	"time"
)

const (
	maxDelay = 2 * time.Minute
	base     = 100 * time.Millisecond
)

// Do is a function x^e multiplied by a factor of 0.1s.
// Result is limited to 2 minute.
// attempts	Pow(attempts, e)	结果（退避时间）
// 1	1^e = 1	100ms
// 2	2^e ≈ 6.59	≈ 659ms
// 3	3^e ≈ 17.44	≈ 1.744s
// 5	5^e ≈ 80.45	≈ 8.045s
// 10	10^e ≈ 1995.26	≈ 199.5s → capped to 2m
// 13	13^e ≈ 3701.28	≈ 370s → capped to 2m
func Do(attempts int) time.Duration {
	if attempts > 13 {
		return maxDelay
	}
	return time.Duration(math.Pow(float64(attempts), math.E)) * base
}

// DoWithJitter returns an exponential backoff duration with full jitter.
// Max backoff time is capped at 2 minutes.
func DoWithJitter(attempts int) time.Duration {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	f := r.Float64() * float64(Do(attempts))
	return time.Duration(f)
}
