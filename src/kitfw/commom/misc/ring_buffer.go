package misc

import (
	"fmt"
	"math"
)

// 一个环形buffer
type RingBuffer struct {
	buffer   []interface{}
	begIndex int
	endIndex int
	max      int
	strict   bool
}

// 生成一个新的环形buffer
// origin代表初始的buffer大小，合适的origin参数可以防止不断生成新的slice
// max代表最大允许存在的buffer大小，如果max<=0则表示不限制，如果超过max限制
// 之后所有的数据都会被丢弃，如果是无限制的情况下，超过maxint32则会panic
// strict如果为true，则表现得是一个严格意义上的fifo，否则则不保证
func NewRingBuffer(origin int, max int, strict bool) *RingBuffer {
	if max <= 0 {
		max = math.MaxInt32
	}

	if origin <= 0 {
		origin = 128
	} else {
		size := 1
		for size < origin {
			size <<= 1
		}
		origin = size
	}

	if origin > max {
		origin = max
	}

	return &RingBuffer{
		buffer: make([]interface{}, origin),
		max:    max,
		strict: strict,
	}
}

func (this *RingBuffer) Push(data interface{}) {
	tsize := len(this.buffer)
	if tsize == this.Size() {
		if tsize == this.max {
			if this.max == math.MaxInt32 {
				panic(fmt.Errorf("buffer overflow"))
			}
			return
		}

		if this.strict {
			nsize := tsize * 2
			if nsize > this.max {
				nsize = this.max
			}
			buffer := make([]interface{}, nsize)
			copy(buffer[0:tsize-this.begIndex], this.buffer[this.begIndex:tsize])
			if this.endIndex > tsize {
				copy(buffer[tsize-this.begIndex:tsize], this.buffer[0:this.endIndex%tsize])
			}
			buffer[tsize] = data
			this.buffer = buffer
		} else {
			this.buffer = append(this.buffer, data)
		}
		this.begIndex = 0
		this.endIndex = tsize + 1
	} else {
		this.buffer[this.endIndex%tsize] = data
		this.endIndex++
	}
}

func (this *RingBuffer) Peek() interface{} {
	if this.Size() == 0 {
		panic(fmt.Errorf("no data"))
	}

	return this.buffer[this.begIndex]
}

func (this *RingBuffer) Pop() interface{} {
	if this.Size() == 0 {
		panic(fmt.Errorf("no data"))
	}

	data := this.buffer[this.begIndex]
	this.begIndex++
	tsize := len(this.buffer)
	if this.begIndex >= tsize {
		this.begIndex -= tsize
		this.endIndex -= tsize
	}

	return data
}

func (this *RingBuffer) Size() int {
	return this.endIndex - this.begIndex
}

/* vim: set ts=4 sw=4 sts=4 tw=100 noet: */
