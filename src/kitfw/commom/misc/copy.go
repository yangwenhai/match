package misc

import (
	"reflect"
)

func Copy(s reflect.Value, t reflect.Value) (err error) {
	tv, err := Clone(s)
	if err != nil {
		return err
	}

	t.Set(tv)
	return nil
}

/*
func vcopy(s reflect.Value, t reflect.Value, level int, m map[uintptr]reflect.Value) error {
	fmt.Println(s, t, t.CanSet())
	switch s.Kind() {
	case reflect.Interface:
		if s.IsNil() {
			t.Set(s)
			return nil
		}
	case reflect.Slice, reflect.Ptr, reflect.Map:
		if s.IsNil() {
			t.Set(s)
			return nil
		}

		if level > 0 {
			addr := s.Pointer()
			v, ok := m[addr]
			if ok {
				t.Set(v)
				return nil
			}
		}
	}

	switch s.Kind() {
	case reflect.Bool, reflect.Float32, reflect.Float64,
		reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Int8, reflect.String, reflect.Uint, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uint8:
		t.Set(s)
	case reflect.Array:
		for i := 0; i < s.Len(); i++ {
			err := vcopy(s.Index(i), t.Index(i), level+1, m)
			if err != nil {
				return err
			}
		}
	case reflect.Slice:
		if t.Len() < s.Len() {
			i := reflect.MakeSlice(s.Type(), s.Len()-t.Len(), s.Len()-t.Len())
			t.Set(reflect.AppendSlice(t, i)) // t= append(t, slice)
		} else if t.Len() > s.Len() {
			t.Set(t.Slice(0, s.Len())) // t=t[0:len]
		}
		m[s.Pointer()] = t
		for i := 0; i < s.Len(); i++ {
			err := vcopy(s.Index(i), t.Index(i), level+1, m)
			if err != nil {
				return err
			}
		}
	case reflect.Map:
		if t.IsNil() {
			t.Set(reflect.MakeMap(t.Type()))
		}
		if t.Len() != 0 {
			t.Set(reflect.MakeMap(s.Type()))
		}
		m[s.Pointer()] = t
		for _, k := range s.MapKeys() {
			sv := s.MapIndex(k)
			tv := reflect.New(sv.Type()).Elem()
			err := vcopy(sv, tv, level+1, m)
			if err != nil {
				return err
			}
			t.SetMapIndex(k, tv)
		}
	case reflect.Struct:
		for i := 0; i < s.NumField(); i++ {
			err := vcopy(s.Field(i), t.Field(i), level+1, m)
			if err != nil {
				return err
			}
		}
	case reflect.Ptr:
		if t.IsNil() {
			t.Set(reflect.New(s.Elem().Type()))
		}
		m[s.Pointer()] = t
		err := vcopy(s.Elem(), t.Elem(), level+1, m)
		if err != nil {
			return err
		}
	case reflect.Interface:
		if t.Elem().IsValid() && t.Elem().Type() == s.Elem().Type() {
			return vcopy(s.Elem(), t.Elem(), level+1, m)
		} else {
			s = s.Elem()
			switch s.Kind() {
			case reflect.Map:
				t.Set(reflect.MakeMap(s.Type()))
				t = t.Elem()
			case reflect.Slice:
				t.Set(reflect.MakeSlice(s.Type(), s.Len(), s.Cap()))
				t = t.Elem()
			case reflect.Ptr:
				t.Set(reflect.New(s.Elem().Type()))
				t = t.Elem()
			case reflect.Struct, reflect.Array:
				tv := reflect.New(s.Type())
				err := vcopy(s, tv.Elem(), level+1, m)
				if err != nil {
					return err
				}
				t.Set(tv.Elem())
				return nil
			}
			return vcopy(s, t, level+1, m)
		}
	default:
		return fmt.Errorf("unsupported type:%s", s.Kind())
	}
	return nil
}
*/

/* vim: set ts=4 sw=4 sts=4 tw=100 noet: */
