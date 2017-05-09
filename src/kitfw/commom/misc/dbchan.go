package misc

import (
	"fmt"
)

// 一个支持动态长度的channel类型
// 利用两个channel以及一个circular buffer来实现
type TDBChan struct {
	in     chan interface{}   // 用于读出数据的channel
	out    chan interface{}   // 用于写入数据的channel
	In     chan<- interface{} // 给用户使用的channel，写入数据
	Out    <-chan interface{} // 给用户使用的channel，写出数据
	buffer *RingBuffer        // 存储中间数据的buffer
}

func (this *TDBChan) close() {
	close(this.in)
	if this.in != this.out {
		close(this.out)
	}
}

// 后台goroutine，用于处理以下事务
// 1 将InChan中写入的数据读出并放入buffer中
// 2 将buffer中的数据读出并放入OutChan中
func (this *TDBChan) start() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf("uncaught exception:%v for DBChan", err)
			this.close()
		}
	}()
	var data, newData interface{}
	for {
		if this.buffer.Size() == 0 {
			data = <-this.in
			if data == nil {
				this.close()
				return
			}
			this.buffer.Push(data)
		} else {
			data = this.buffer.Peek()
			select {
			case this.out <- data:
				this.buffer.Pop()
			case newData = <-this.in:
				if newData == nil {
					return
				}
				this.buffer.Push(newData)
			}
		}
	}
}

func (this *TDBChan) Close() {
	this.In <- nil
}

func NewDBChan(origin int, max int, strict bool) *TDBChan {
	c := &TDBChan{
		buffer: NewRingBuffer(origin, max, strict),
		in:     make(chan interface{}),
	}

	if strict {
		c.out = make(chan interface{})
	} else {
		c.out = c.in
	}

	c.In = c.in
	c.Out = c.out
	go c.start()
	return c
}

/* vim: set ts=4 sw=4 sts=4 tw=100 noet: */
