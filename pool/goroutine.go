// @Author Eric
// @Date 2024/6/4 13:13:00
// @Desc
package pool

import (
	"github.com/panjf2000/ants/v2"
	"time"
)

type Goroutine struct {
	pool *ants.Pool
}

func NewGoroutine(size int) *Goroutine {
	if size <= 0 {
		size = 2000
	}
	p, err := ants.NewPool(size, ants.WithExpiryDuration(10*time.Second))
	if err != nil {
		panic(err)
	}
	pool := &Goroutine{
		pool: p,
	}
	return pool
}
