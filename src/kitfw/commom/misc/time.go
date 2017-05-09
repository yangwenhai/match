package misc

import (
	"sync"
	"sync/atomic"
	"time"
)

var (
	rwLock sync.RWMutex
	ticker *iTicker
)

type iTicker struct {
	ticker   *time.Ticker
	unixNano int64
	d        time.Duration
}

func (this *iTicker) start() {
	go this.handleTicker()
}

func (this *iTicker) handleTicker() {
	for {
		newT, ok := <-this.ticker.C
		if !ok {
			break
		}
		atomic.StoreInt64(&this.unixNano, newT.UnixNano())
	}
}

func (this *iTicker) now() time.Time {
	nano := atomic.LoadInt64(&this.unixNano)
	return time.Unix(0, nano)
}

func (this *iTicker) stop() {
	this.ticker.Stop()
}

func newTicker(d time.Duration) *iTicker {
	ret := &iTicker{
		ticker:   time.NewTicker(d),
		unixNano: time.Now().UnixNano(),
		d:        d,
	}

	return ret
}

// 得到当前的时间，这个时间是以d指定的精度来进行更新的
func Now(d time.Duration) time.Time {
	rwLock.RLock()

	if ticker == nil || ticker.d > d {
		rwLock.RUnlock()
		rwLock.Lock()
		defer rwLock.Unlock()
		if ticker != nil {
			ticker.stop()
		}
		ticker = newTicker(d)
		ticker.start()
	} else {
		defer rwLock.RUnlock()
	}

	return ticker.now()
}

// 如果当前正在运行的ticker精度比d指定的还小则删除当前的ticker
// 否则不做任何动作
func DestroyTicker(d time.Duration) {
	rwLock.RLock()
	defer rwLock.RUnlock()

	if ticker == nil || ticker.d >= d {
		return
	}

	rwLock.Lock()
	ticker.stop()
	ticker = nil
	rwLock.Unlock()
}

/* vim: set ts=4 sw=4 sts=4 tw=100 noet: */
