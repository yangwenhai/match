// 一些常用的方法
package misc

import (
	"fmt"
	"reflect"
	"unicode"
)

// 检查一个名字是否是可以被外部包访问的
func Exportable(name string) bool {
	if unicode.IsUpper([]rune(name)[0]) {
		return true
	}

	return false
}

func Indirect(v reflect.Value) (reflect.Value, error) {
	if !v.CanSet() {
		return v, NewError("value type:%s can't be set", v.Type().String())
	}

	for {
		if nilable(v) && v.IsNil() {
		}
	}

}

func initValue(typ reflect.Type, size, cap int) reflect.Value {
	switch typ.Kind() {
	case reflect.Chan:
		return reflect.MakeChan(typ, size)
	case reflect.Map:
		return reflect.MakeMap(typ)
	case reflect.Func:
		return reflect.Zero(typ)
	case reflect.Slice:
		return reflect.MakeSlice(typ, size, cap)
	case reflect.Ptr:
		return reflect.New(typ.Elem())
	case reflect.Interface:
	}
	return reflect.Value{}
}

// 判断一个类型是否可以调用Nil方法
func nilable(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan, reflect.Map, reflect.Func, reflect.Ptr, reflect.Interface, reflect.Slice:
		return true
	default:
		return false
	}
}

// 判断一个类型是否需要调用对应的make方法
func makeable(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan, reflect.Map, reflect.Func, reflect.Slice:
		return true
	default:
		return false
	}
}

// 深度clone方法，只支持slice,map以及pointer类型
// 以上三种类型以地址作为判断是否相同的唯一条件
func Clone(s reflect.Value) (reflect.Value, error) {
	switch s.Kind() {
	case reflect.Slice, reflect.Ptr, reflect.Map, reflect.Struct:
	default:
		return s, fmt.Errorf("clone type must be slice/pointer/map")
	}
	m := make(map[uintptr]reflect.Value)
	return clone(s, m)
}

func clone(s reflect.Value, m map[uintptr]reflect.Value) (reflect.Value, error) {
	switch s.Kind() {
	case reflect.Interface:
		if s.IsNil() {
			return s, nil
		}
	case reflect.Slice, reflect.Ptr, reflect.Map:
		if s.IsNil() {
			return s, nil
		}

		addr := s.Pointer()
		v, ok := m[addr]
		if ok {
			return v, nil
		}
	}

	switch s.Kind() {
	case reflect.Array, reflect.Bool, reflect.Float32, reflect.Float64,
		reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Int8, reflect.String, reflect.Uint, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uint8:
		return s, nil
	case reflect.Slice:
		t := reflect.MakeSlice(s.Type(), s.Len(), s.Cap())
		m[s.Pointer()] = t
		for i := 0; i < s.Len(); i++ {
			elem, err := clone(s.Index(i), m)
			if err != nil {
				return s, err
			}
			t.Index(i).Set(elem)
		}
		return t, nil
	case reflect.Map:
		t := reflect.MakeMap(s.Type())
		m[s.Pointer()] = t
		for _, k := range s.MapKeys() {
			v, err := clone(s.MapIndex(k), m)
			if err != nil {
				return s, err
			}
			t.SetMapIndex(k, v)
		}
		return t, nil
	case reflect.Struct:
		t := reflect.New(s.Type()).Elem()
		for i := 0; i < s.NumField(); i++ {
			f, err := clone(s.Field(i), m)
			if err != nil {
				return s, err
			}
			t.Field(i).Set(f)
		}
		return t, nil
	case reflect.Ptr:
		t := reflect.New(s.Elem().Type())
		m[s.Pointer()] = t
		v, err := clone(s.Elem(), m)
		if err != nil {
			return s, err
		}
		t.Elem().Set(v)
		return t, nil
	case reflect.Interface:
		return clone(s.Elem(), m)
	default:
		return s, fmt.Errorf("unsupported type:%v", s)
	}

}

// 对所有不希望panic的go调用函数加一个壳
// 如果希望由编译器进行参数检查，则可以直接使用
// go func(){
// 	defer ...
// 	call ...
// }()
func SafeGo(data ...interface{}) {
	fn := reflect.ValueOf(data)
	if fn.Elem().Kind() != reflect.Func || fn.Type().NumIn() != len(data)-1 {
		panic(fmt.Errorf("invalid call for func:%v", data))
	}

	var in []reflect.Value
	for _, arg := range data[1:] {
		in = append(in, reflect.ValueOf(arg))
	}

	go func() {
		defer func() {
			err := recover()
			fmt.Printf("uncaught err:%v", err)
		}()
		fn.CallSlice(in)
	}()
}

/* vim: set ts=4 sw=4 sts=4 tw=100 noet: */
