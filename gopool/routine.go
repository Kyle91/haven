package gopool

import (
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

type Goroutine struct {
	pool *ants.Pool
}

var instance *Goroutine
var once sync.Once

func GetInstance(size int, expiryDuration time.Duration) *Goroutine {
	once.Do(func() {
		if size <= 0 {
			size = 2000
		}
		if expiryDuration <= 0 {
			expiryDuration = 10 * time.Second
		}
		p, err := ants.NewPool(size, ants.WithExpiryDuration(expiryDuration))
		if err != nil {
			panic(err)
		}
		instance = &Goroutine{
			pool: p,
		}
	})
	return instance
}

func (g *Goroutine) Submit(task func()) error {
	return g.pool.Submit(task)
}

func (g *Goroutine) Release() {
	g.pool.Release()
}
