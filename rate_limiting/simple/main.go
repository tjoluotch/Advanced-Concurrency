package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Limiter struct {
	mu         sync.Mutex
	rate       int
	bucketSize int
	nTokens    int
	// lastTokenTime represents the last token time
	lastTokenTime time.Time
}

func (lt *Limiter) Wait() {
	lt.mu.Lock()
	defer lt.mu.Unlock()

	if lt.nTokens > 0 {
		lt.nTokens--
		return
	}

	tElapsed := time.Since(lt.lastTokenTime)
	period := time.Second / time.Duration(lt.rate)
	nTokens := tElapsed.Nanoseconds() / period.Nanoseconds()
	lt.nTokens = int(nTokens)
	if lt.nTokens > lt.bucketSize {
		lt.nTokens = lt.bucketSize
	}
	lt.lastTokenTime = lt.lastTokenTime.Add(time.Duration(nTokens) * period)

	if lt.nTokens > 0 {
		lt.nTokens--
		return
	}

	next := lt.lastTokenTime.Add(period)
	wait := next.Sub(time.Now())

	if wait >= 0 {
		time.Sleep(wait)
	}
	lt.lastTokenTime = next
}

func NewLimiter(rate, limit int) *Limiter {
	return &Limiter{
		mu:            sync.Mutex{},
		rate:          rate,
		bucketSize:    limit,
		nTokens:       limit,
		lastTokenTime: time.Now(),
	}
}

func main() {
	limiter := NewLimiter(5, 10)

	for i := 0; i < 100; i++ {
		limiter.Wait()
		fmt.Printf("Request %v %+v\n", time.Now(), limiter)
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(400)))

	}
	time.Sleep(time.Second * 2)
	for i := 0; i < 100; i++ {
		limiter.Wait()
		fmt.Printf("Request %v %+v\n", time.Now(), limiter)
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(400)))
	}
}
