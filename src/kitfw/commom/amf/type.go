package amf

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

const (
	UNDEFINED_MARKER = 0x00
	NULL_MARKER      = 0x01
	FALSE_MARKER     = 0x02
	TRUE_MARKER      = 0x03
	INTEGER_MARKER   = 0x04
	DOUBLE_MARKER    = 0x05
	STRING_MARKER    = 0x06
	XMLDOC_MARKER    = 0x07
	DATE_MARKER      = 0x08
	ARRAY_MARKER     = 0x09
	OBJECT_MARKER    = 0x0a
	XML_MARKER       = 0x0b
	BYTEARRAY_MARKER = 0x0c
)

//编码器的对象
type Encoder struct {
	//编码后所写入的字节流
	writer io.Writer
	//当前编码所处的位置，用于出错时提示使用
	KeyStack
	//字符串的缓存池
	scache map[string]int
	//对象的缓存池
	ocache map[uintptr]int
	//当前的字节序
	byteOrder binary.ByteOrder
}

//解码器的对象
type Decoder struct {
	//解码时用于读取输入的字节流
	reader io.Reader
	//当前解码所处的位置
	KeyStack
	//字符串的缓存池
	scache [][]byte
	//对象的缓存池
	ocache []reflect.Value
	//当前读了多少字节
	bytes int
	//当前的字符序
	byteOrder binary.ByteOrder
}

//解码错误
type AmfError struct {
	message string
	stack   string
}

type KeyStack []string

func (err *AmfError) Error() string {
	return fmt.Sprintf("err:%s, stack:%s", err.message, err.stack)
}

func NewAmfError(stack KeyStack, message string) *AmfError {
	return &AmfError{
		stack:   stack.Stack(),
		message: message,
	}
}

func (stack *KeyStack) pushKey(key string) {
	*stack = append(*stack, key)
}

func (stack *KeyStack) popKey() string {
	length := len(*stack)
	key := ""
	if length > 0 {
		key = (*stack)[length-1]
		*stack = (*stack)[0 : length-1]
	}
	return key
}

func (stack *KeyStack) Stack() string {
	ret := ""
	for _, key := range *stack {
		key = "{" + key + "}"
		if ret == "" {
			ret = key
		} else {
			ret = ret + "->" + key
		}
	}
	return ret
}

type AmfNamer interface {
	GetFieldName(amfName string) string
	GetAmfName(fieldName string) string
}

/* vim: set ts=4 sw=4 sts=4 tw=100 noet: */
