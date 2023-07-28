package ticker

import (
	"time"
)

type (
	Token struct {
		fn    func() string
		value string
		ttl   time.Duration
	}
)

func NewToken(fn func() string, ttl time.Duration) *Token {
	return &Token{
		fn:  fn,
		ttl: ttl,
	}
}

func (t *Token) Get() string {
	if t.value != "" {
		return t.value
	}

	<-t.retrieve()
	return t.value
}

func (t *Token) retrieve() <-chan struct{} {
	done := make(chan struct{})
	go func() {
		for {
			t.value = t.fn()
			done <- struct{}{}
			time.Sleep(t.ttl)
		}
	}()

	return done
}
