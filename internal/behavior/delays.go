package behavior

import (
	"math/rand"
	"time"
)

// SleepRange sleeps for a random duration between min and max
func SleepRange(min, max time.Duration) {
	if max < min {
		max = min
	}
	duration := min + time.Duration(rand.Int63n(int64(max-min+1)))
	time.Sleep(duration)
}

// Jitter adds random variance to a duration
func Jitter(d time.Duration, pct float64) time.Duration {
	if pct <= 0 {
		return d
	}
	variance := float64(d) * pct
	offset := (rand.Float64()*2 - 1) * variance
	return d + time.Duration(offset)
}

// Backoff calculates exponential backoff with max limit
func Backoff(base time.Duration, factor float64, max time.Duration, attempt int) time.Duration {
	duration := base
	for i := 0; i < attempt; i++ {
		duration = time.Duration(float64(duration) * factor)
		if duration > max {
			return max
		}
	}
	return duration
}

// RandomPause adds a random human-like pause
func RandomPause() {
	SleepRange(100*time.Millisecond, 500*time.Millisecond)
}

// LongPause adds a longer random pause
func LongPause() {
	SleepRange(1*time.Second, 3*time.Second)
}
