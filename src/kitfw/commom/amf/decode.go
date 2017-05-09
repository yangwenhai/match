package amf

import (
	"encoding/binary"
	"io"
	"reflect"
	"strconv"
)

func NewDecoder(reader io.Reader, byteOrder binary.ByteOrder) *Decoder {
	return &Decoder{
		reader:    reader,
		byteOrder: byteOrder,
		bytes:     0,
	}
}

//解码，要求data只能以指针形式提供
func (decoder *Decoder) Decode(data interface{}) (err error) {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return NewAmfError(decoder.KeyStack, "non-nil pointer expected")
	}

	decoder.Reset()
	return decoder.readValue(v)
}

//将当前的解码所有参数重置
func (decoder *Decoder) Reset() {
	decoder.ocache = nil
	decoder.scache = nil
	decoder.KeyStack = nil
	decoder.bytes = 0
}

func (decoder *Decoder) readStringBytes() ([]byte, error) {
	index, err := decoder.readU29()
	if err != nil {
		return nil, err
	}

	var ret []byte
	if (index & 0x01) == 0 {
		ret = decoder.scache[int(index>>1)]
	} else {
		index >>= 1
		ret, err = decoder.readBytes(int(index))
		if err != nil {
			return nil, err
		}

		if len(ret) != 0 {
			decoder.scache = append(decoder.scache, ret)
		}
	}

	return ret, nil
}

func (decoder *Decoder) readString(value reflect.Value, binary bool) error {
	ret, err := decoder.readStringBytes()
	if err != nil {
		return err
	}
	value = decoder.indirect(value, false)
	return decoder.setString(value, ret, binary)
}

func (decoder *Decoder) setString(value reflect.Value, ret []byte, binary bool) error {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		num, err := strconv.ParseInt(string(ret), 0, 64)
		if err != nil {
			return err
		}
		value.SetInt(num)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		num, err := strconv.ParseUint(string(ret), 0, 64)
		if err != nil {
			return err
		}
		value.SetUint(num)
	case reflect.Slice:
		if value.Type().Elem().Kind() != reflect.Int8 {
			return NewAmfError(decoder.KeyStack, "slice type not int8 for string")
		}
		value.SetBytes(ret)
	case reflect.Array:
		if value.Type().Elem().Kind() != reflect.Int8 {
			return NewAmfError(decoder.KeyStack, "array type not int8 for string")
		}
		if value.Type().Len() < len(ret) {
			return NewAmfError(decoder.KeyStack, "array length too short for string")
		}
		reflect.Copy(value, reflect.ValueOf(ret))
	case reflect.String:
		value.SetString(string(ret))
	case reflect.Interface:
		if binary {
			value.Set(reflect.ValueOf(ret))
		} else {
			value.Set(reflect.ValueOf(string(ret)))
		}
	default:
		return NewAmfError(decoder.KeyStack, "unsupported type for string")
	}
	return nil
}

func (decoder *Decoder) readBool(value reflect.Value, b bool) error {
	value = decoder.indirect(value, false)
	switch value.Kind() {
	case reflect.Bool:
		value.SetBool(b)
	case reflect.Interface:
		value.Set(reflect.ValueOf(b))
	case reflect.String:
		value.SetString(strconv.FormatBool(b))
	default:
		return NewAmfError(decoder.KeyStack, "unsupported type for bool")
	}
	return nil
}

func (decoder *Decoder) readNil(value reflect.Value) error {
	value = decoder.indirect(value, false)
	switch value.Kind() {
	case reflect.Interface, reflect.Ptr, reflect.Map, reflect.Slice:
		if !value.IsNil() {
			value.Set(reflect.ValueOf(nil))
		}
	default:
		return NewAmfError(decoder.KeyStack, "unsupported type for nil")
	}
	return nil
}

func (decoder *Decoder) readU29() (uint32, error) {
	ret := uint32(0)
	for i := 0; i < 4; i++ {
		b, err := decoder.readByte()
		if err != nil {
			return 0, err
		}

		if i != 3 {
			ret = (ret << 7) | uint32(b&0x7f)
			if (b & 0x80) == 0 {
				break
			}
		} else {
			ret = (ret << 8) | uint32(b)
		}
	}

	return ret, nil
}

func (decoder *Decoder) readInteger(value reflect.Value) error {

	uv, err := decoder.readU29()
	if err != nil {
		return err
	}

	vv := int32(uv)
	if uv > 0xfffffff {
		vv = int32(uv - 0x20000000)
	}

	value = decoder.indirect(value, false)
	switch value.Kind() {
	case reflect.Interface:
		value.Set(reflect.ValueOf(vv))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value.SetInt(int64(vv))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value.SetUint(uint64(uv))
	case reflect.String:
		value.SetString(strconv.FormatInt(int64(vv), 10))
	default:
		return NewAmfError(decoder.KeyStack, "unsupported type for integer")
	}
	return nil
}

func (decoder *Decoder) readArray(value reflect.Value) error {
	index, err := decoder.readU29()
	if err != nil {
		return err
	}

	value = decoder.indirect(value, false)

	if (index & 0x01) == 0 {
		slice := decoder.ocache[int(index>>1)]
		value.Set(slice)
		return nil
	}

	index >>= 1
	sep, err := decoder.readByte()
	if err != nil {
		return err
	}

	if sep != 0x01 {
		return NewAmfError(decoder.KeyStack, "ecma array not allowed")
	}

	decoder.ocache = append(decoder.ocache, value)

	switch value.Kind() {
	case reflect.Interface:
		var slice []interface{}
		value.Set(reflect.MakeSlice(reflect.TypeOf(slice), int(index), int(index)))
		value = value.Elem()
	case reflect.Array:
	case reflect.Slice:
		if value.IsNil() {
			value.Set(reflect.MakeSlice(value.Type(), int(index), int(index)))
		}
	//php中空map (array)转换为amf3中的array类型
	case reflect.Map:
		if index == 0 {
			break
		}
		fallthrough
	default:
		return NewAmfError(decoder.KeyStack, "unsupported type for array")
	}

	if value.Len() < int(index) {
		return NewAmfError(decoder.KeyStack, "length not enough")
	}

	for i := 0; i < int(index); i++ {
		decoder.pushKey(strconv.Itoa(i))
		c := value.Index(i)
		err = decoder.readValue(c)
		if err != nil {
			return err
		}
		decoder.popKey()
	}

	return nil
}

func (decoder *Decoder) readObject(value reflect.Value) error {
	index, err := decoder.readU29()
	if err != nil {
		return err
	}

	value = decoder.indirect(value, false)

	if (index & 0x01) == 0 {
		value.Set(decoder.ocache[int(index>>1)])
		return nil
	}

	if index != 0x0b {
		return NewAmfError(decoder.KeyStack, "invalid object type")
	}

	sep, err := decoder.readByte()
	if err != nil {
		return err
	}

	if sep != 0x01 {
		return NewAmfError(decoder.KeyStack, "type object not allowed")
	}

	decoder.ocache = append(decoder.ocache, value)

	var namer AmfNamer
	ok := false
	switch value.Kind() {
	case reflect.Interface:
		var m map[string]interface{}
		mtype := reflect.TypeOf(m)
		value.Set(reflect.MakeMap(mtype))
		value = value.Elem()
	case reflect.Map:
		if value.IsNil() {
			value.Set(reflect.MakeMap(value.Type()))
		}
	case reflect.Struct:
		namer, ok = value.Addr().Interface().(AmfNamer)
	default:
		return NewAmfError(decoder.KeyStack, "unsupported type for object")
	}

	for {
		keyBytes, err := decoder.readStringBytes()
		if err != nil {
			return err
		}

		key := string(keyBytes)

		if key == "" {
			break
		}

		decoder.pushKey(key)

		var v reflect.Value

		switch {
		case value.Kind() == reflect.Map:
			v = reflect.New(value.Type().Elem())
		case value.Kind() == reflect.Struct:
			if ok {
				key = namer.GetFieldName(key)
			} else {
				t := value.Type()
				for i := 0; i < t.NumField(); i++ {
					f := t.Field(i)
					if key == f.Tag.Get("amf") {
						key = f.Name
						break
					}
				}
			}
			v = value.FieldByName(key)
			if err != nil {
				return err
			}
		}

		err = decoder.readValue(v)
		if err != nil {
			return err
		}

		/*
			switch {
			case value.Kind() == reflect.Map:
				value.SetMapIndex(reflect.ValueOf(key), v.Elem())
			}
		*/

		switch {
		case value.Kind() == reflect.Map:
			keyKind := value.Type().Key().Kind()
			if keyKind == reflect.String {
				value.SetMapIndex(reflect.ValueOf(key), v.Elem())
			} else {
				switch keyKind {
				case reflect.Int, reflect.Int64:
					num, err := strconv.ParseInt(key, 0, 64)
					if err != nil {
						return err
					}
					value.SetMapIndex(reflect.ValueOf(int(num)), v.Elem())
				case reflect.Uint, reflect.Uint64:
					num, err := strconv.ParseUint(key, 0, 64)
					if err != nil {
						return err
					}
					value.SetMapIndex(reflect.ValueOf(uint(num)), v.Elem())
				default:
					return NewAmfError(decoder.KeyStack, "invalid map key type")
				}
			}
		}

		decoder.popKey()
	}

	return nil
}

func (decoder *Decoder) readDouble(value reflect.Value) error {
	value = decoder.indirect(value, false)
	var f float64
	err := binary.Read(decoder.reader, decoder.byteOrder, &f)
	if err != nil {
		return err
	}
	switch value.Kind() {
	case reflect.Interface:
		value.Set(reflect.ValueOf(f))
	case reflect.Float32, reflect.Float64:
		value.SetFloat(f)
	case reflect.String:
		value.SetString(strconv.FormatFloat(f, 'f', 10, 64))
	//amf int只有29位， 大于int29可能被转换为double了
	case reflect.Int, reflect.Int32, reflect.Int64:
		value.SetInt(int64(f))
	default:
		return NewAmfError(decoder.KeyStack, "unsupported type for double")
	}
	return nil
}

func (decoder *Decoder) readValue(value reflect.Value) error {
	if !value.IsValid() {
		return NewAmfError(decoder.KeyStack, "invalid value")
	}

	marker, err := decoder.readByte()
	if err != nil {
		return err
	}

	switch marker {
	case NULL_MARKER:
		return decoder.readNil(value)
	case FALSE_MARKER:
		return decoder.readBool(value, false)
	case TRUE_MARKER:
		return decoder.readBool(value, true)
	case INTEGER_MARKER:
		return decoder.readInteger(value)
	case DOUBLE_MARKER:
		return decoder.readDouble(value)
	case STRING_MARKER:
		return decoder.readString(value, false)
	case XMLDOC_MARKER:
		return decoder.readString(value, true)
	case XML_MARKER:
		return decoder.readString(value, true)
	case BYTEARRAY_MARKER:
		return decoder.readString(value, true)
	case DATE_MARKER:
		return decoder.readDouble(value)
	case ARRAY_MARKER:
		return decoder.readArray(value)
	case OBJECT_MARKER:
		return decoder.readObject(value)
	default:
		return NewAmfError(decoder.KeyStack, "unsupported marker")
	}
}

func (decoder *Decoder) indirect(v reflect.Value, decodingNull bool) reflect.Value {
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		v = v.Addr()
	}
	for {
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() && (!decodingNull || e.Elem().Kind() == reflect.Ptr) {
				v = e
				continue
			}
		}

		if v.Kind() != reflect.Ptr {
			break
		}

		if v.Elem().Kind() != reflect.Ptr && decodingNull && v.CanSet() {
			break
		}

		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}

	return v

}

func (decoder *Decoder) readByte() (byte, error) {
	bytes, err := decoder.readBytes(1)
	if err != nil {
		return 0, err
	}

	return bytes[0], nil
}

func (decoder *Decoder) readBytes(n int) ([]byte, error) {
	if n < 0 {
		return nil, NewAmfError(decoder.KeyStack, "invalid bytes number")
	}

	if n == 0 {
		return []byte{}, nil
	}

	bytes := make([]byte, n)
	nn, err := decoder.reader.Read(bytes)
	if nn != n || err != nil {
		return nil, NewAmfError(decoder.KeyStack, "end of file")
	}

	decoder.bytes += n

	return bytes, nil
}

/* vim: set ts=4 sw=4 sts=4 tw=100 noet: */
