package misc

import (
	"fmt"
	"time"
)

// 一个支持动态长度的channel类型
// 利用两个channel以及一个circular buffer来实现
type DelayChan struct {
	In      chan<- interface{} // 只可以写入的channel
	in      chan interface{}
	Out     <-chan interface{} // 只可以读出的channel
	out     chan interface{}
	ticker  <-chan time.Time // ticker，用于处理时间
	input   *RingBuffer      // 输入的buffer
	output  *RingBuffer
	delay   time.Duration
	running bool
}

// 一个带时间信息的
type timedObject struct {
	t    time.Time
	data interface{}
}

// 后台goroutine，用于处理以下事务
// 1 将InChan中写入的数据读出并放入buffer中
// 2 将buffer中的数据读出并放入OutChan中
func (this *DelayChan) start() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf("uncaught exception:%v for DBChan", err)
			close(this.in)
			close(this.out)
		}
	}()

	this.running = true
	now := time.Now()
	var idata, odata interface{}
	for this.running {
		if this.output.Size() != 0 {
			odata = this.output.Peek()
			select {
			case idata = <-this.in:
				this.input.Push(timedObject{now, idata})
			case now = <-this.ticker:
				if this.input.Size() == 0 {
					continue
				}

				data := this.input.Peek().(timedObject)
				if data.t.Add(this.delay).Before(now) {
					this.input.Pop()
					this.output.Push(data.data)
				}
			case this.out <- odata:
				this.output.Pop()
			}
		} else {
			select {
			case idata = <-this.in:
				this.input.Push(timedObject{now, idata})
			case now = <-this.ticker:
				if this.input.Size() == 0 {
					continue
				}

				data := this.input.Peek().(timedObject)
				if data.t.Add(this.delay).Before(now) {
					this.input.Pop()
					this.output.Push(data.data)
				}
			}
		}
	}
}

func (this *DelayChan) Stop() {
	this.running = false
}

func NewDelayChan(precision time.Duration, delay time.Duration) *DelayChan {
	c := &DelayChan{
		input:  NewRingBuffer(128, 0, true),
		output: NewRingBuffer(128, 0, true),
		in:     make(chan interface{}),
		out:    make(chan interface{}),
		ticker: time.Tick(precision),
		delay:  delay,
	}
	c.In = c.in
	c.Out = c.out
	go c.start()
	return c
}

/* vim: set ts=4 sw=4 sts=4 tw=100 noet: */
