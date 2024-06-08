// @Author Eric
// @Date 2024/6/4 13:13:00
// @Desc
package gopool

import (
	"errors"
	"github.com/panjf2000/ants/v2"
	"time"
)

type Routine struct {
	pool *ants.Pool
}

// NewGoroutine
//
//	@Description: 创建一个新的 Goroutine 池
//	@param size 池大小 默认2000
//	@param expiryDuration 任务过期时间 默认10秒
//	@return *Goroutine
//	@return error
func NewGoroutine(size int, expiryDuration time.Duration) (*Routine, error) {
	if size <= 0 {
		size = 2000
	}
	if expiryDuration <= 0 {
		expiryDuration = 10 * time.Second
	}
	p, err := ants.NewPool(size, ants.WithExpiryDuration(expiryDuration))
	if err != nil {
		return nil, err
	}
	return &Routine{
		pool: p,
	}, nil
}

// Submit 提交任务到 Goroutine 池
func (g *Routine) Submit(task func()) error {
	if g.pool == nil {
		return errors.New("pool is not initialized")
	}
	return g.pool.Submit(func() {
		task()
	})
}

// Release 释放 Goroutine 池
func (g *Routine) Release() {
	if g.pool != nil {
		g.pool.Release()
	}
}
