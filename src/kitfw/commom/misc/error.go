package misc

import (
	"fmt"
	"runtime"
	"strings"
)

// 带有调用栈信息的error接口实现
type StackError struct {
	message string
	stack   [1024]byte
	nbytes  int
}

func (this *StackError) Error() string {
	return this.message
}

// 获取调用栈信息
func (this *StackError) Stack() string {
	return string(this.stack[0:this.nbytes])
}

// 新生成一个*StackError, 错误信息支持printf的format格式
func NewStackError(f string, args ...interface{}) *StackError {
	err := new(StackError)
	err.nbytes = runtime.Stack(err.stack[0:], false)
	err.message = fmt.Sprintf(f, args...)
	err.message = fmt.Sprintf("%s stack:%s", err.message, err.Stack())
	return err
}

var pathPrefix string

func SetPathPrefix(p string) {
	pathPrefix = p
}

func getFileName(file string) string {
	if pathPrefix == "" {
		return file
	}
	return strings.TrimPrefix(file, pathPrefix)
}

func NewError(f string, args ...interface{}) error {
	_, file, num, ok := runtime.Caller(1)
	if ok {
		content := fmt.Sprintf(f, args...)
		file = getFileName(file)
		return fmt.Errorf("%s:%d %s", file, num, content)
	} else {
		return fmt.Errorf(f, args...)
	}
}

/* vim: set ts=4 sw=4 sts=4 tw=100 noet: */
