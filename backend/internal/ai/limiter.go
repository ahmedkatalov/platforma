package ai

import (
	"sync"
	"time"
)

// RateLimiter — простой оконный лимитер «не больше N событий за window» на ключ
// (обычно ключ = id пользователя). Защищает от накрутки и от лишних трат.
type RateLimiter struct {
	mu     sync.Mutex
	window time.Duration
	limit  int
	hits   map[string][]time.Time
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		window: window,
		limit:  limit,
		hits:   make(map[string][]time.Time),
	}
}

// Allow регистрирует попытку и возвращает false, если лимит на окно исчерпан.
func (l *RateLimiter) Allow(key string) bool {
	now := time.Now()
	cutoff := now.Add(-l.window)

	l.mu.Lock()
	defer l.mu.Unlock()

	kept := l.hits[key][:0]
	for _, t := range l.hits[key] {
		if t.After(cutoff) {
			kept = append(kept, t)
		}
	}

	if len(kept) >= l.limit {
		l.hits[key] = kept
		return false
	}
	l.hits[key] = append(kept, now)

	// Ленивая уборка: чтобы карта не росла бесконечно от разовых ключей.
	if len(l.hits) > 10000 {
		for k, ts := range l.hits {
			if len(ts) == 0 || ts[len(ts)-1].Before(cutoff) {
				delete(l.hits, k)
			}
		}
	}
	return true
}
