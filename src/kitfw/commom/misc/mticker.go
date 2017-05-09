package misc

import (
	"sync"
	"time"
)

// 一个允许注册多个channel的ticker，之所以这样做是为了减少系统调用
type MTicker struct {
	ticker   *time.Ticker
	notifier map[string]chan<- time.Time
	lock     sync.Mutex
}

// 生成一个ticker
func NewMTicker(d time.Duration) *MTicker {
	return &MTicker{
		ticker:   time.NewTicker(d),
		notifier: make(map[string]chan<- time.Time),
	}
}

// 开始当前这个ticker
func (this *MTicker) Start() {
	go this.proc()
}

// 处理事件
func (this *MTicker) proc() {
	for {
		select {
		case t := <-this.ticker.C:
			this.broadcast(t)
		}
	}
}

// 广播消息
func (this *MTicker) broadcast(t time.Time) {
	this.lock.Lock()
	defer this.lock.Unlock()

	for _, c := range this.notifier {
		select {
		case c <- t:
		default:
		}
	}
}

// 停止ticker
func (this *MTicker) Stop() {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.ticker.Stop()
	t := time.Time{}
	for _, c := range this.notifier {
		c <- t
	}
}

// 将一个channel绑定到该ticker
func (this *MTicker) Attach(name string, c chan<- time.Time) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	_, ok := this.notifier[name]
	if ok {
		return false
	}

	this.notifier[name] = c
	return true
}

// 将一个channel从这里删除
func (this *MTicker) Detach(name string) {
	this.lock.Lock()
	defer this.lock.Unlock()

	c, ok := this.notifier[name]
	if !ok {
		return
	}

	//c <- time.Time{}
	close(c)
	delete(this.notifier, name)
}

/* vim: set ts=4 sw=4 sts=4 tw=100 noet: */
