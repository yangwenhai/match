package misc

import (
	"fmt"
	"reflect"
	"testing"
)

type T struct {
	Int       int
	String    string
	Map       map[string]int
	Slice     []int
	Array     [10]*string
	Pointer   *T
	Interface interface{}
}

func TestCopy(ti *testing.T) {
	var s T
	s.Int = 10
	s.String = "hello"
	s.Array[0] = &s.String
	s.Map = map[string]int{"Hello": 123}
	s.Slice = []int{1, 2, 3}
	s.Pointer = &s
	s.Interface = s
	var t, sraw interface{}
	t = T{}
	sraw = s
	err := Copy(reflect.ValueOf(&sraw), reflect.ValueOf(&t))
	if err != nil {
		ti.Fatal(err)
		ti.Fail()
	}
	fmt.Println(s)
	fmt.Println(t)
}

/* vim: set ts=4 sw=4 sts=4 tw=100 noet: */
