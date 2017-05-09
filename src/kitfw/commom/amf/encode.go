package amf

import (
	"encoding/binary"
	"fmt"
	"io"
	"kitfw/commom/misc"
	"reflect"
	"strconv"
)

func (encoder *Encoder) Reset() {
	encoder.ocache = make(map[uintptr]int)
	encoder.scache = make(map[string]int)
	encoder.KeyStack = nil
}

func (encoder *Encoder) writeBool(value bool) error {

	if value {
		return encoder.writeByte(TRUE_MARKER)
	} else {
		return encoder.writeByte(FALSE_MARKER)
	}
}

func (encoder *Encoder) writeNull() error {

	return encoder.writeByte(NULL_MARKER)
}

func (encoder *Encoder) writeUint(value uint64) error {

	if value >= 0xfffffff {
		return encoder.writeString(strconv.FormatUint(value, 10))
	}

	err := encoder.writeByte(INTEGER_MARKER)
	if err != nil {
		return err
	}

	return encoder.writeU29(uint32(value))
}

func (encoder *Encoder) writeInt(value int64) error {
	if value > 0xfffffff || value <= -0x10000000 {
		return encoder.writeString(strconv.FormatInt(value, 10))
	}

	err := encoder.writeByte(INTEGER_MARKER)
	if err != nil {
		return err
	}

	return encoder.writeU29(uint32(value & 0x1fffffff))
}

func (encoder *Encoder) writeFloat(value float64) error {

	encoder.writeByte(DOUBLE_MARKER)
	return binary.Write(encoder.writer, encoder.byteOrder, value)
}

func (encoder *Encoder) writeString(value string) error {

	err := encoder.writeByte(STRING_MARKER)
	if err != nil {
		return err
	}

	return encoder.setString(value)
}

func (encoder *Encoder) writeMap(value reflect.Value) error {

	err := encoder.writeByte(OBJECT_MARKER)
	if err != nil {
		return err
	}

	index, ok := encoder.ocache[value.Pointer()]
	if ok {
		index <<= 1
		encoder.writeU29(uint32(index << 1))
		return nil
	}

	err = encoder.writeByte(0x0b)
	if err != nil {
		return err
	}

	err = encoder.setString("")
	if err != nil {
		return err
	}

	keys := value.MapKeys()
	for i := 0; i < len(keys); i++ {
		key := keys[i]
		encoder.pushKey(key.String())

		str := ""
		switch key.Kind() {
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			str = fmt.Sprintf("%d", key.Uint())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			str = fmt.Sprintf("%d", key.Int())
		case reflect.String:
			str = key.String()
		default:
			return NewAmfError(encoder.KeyStack, "unsupported type:"+key.Type().String())
		}

		err = encoder.setString(str)

		if err != nil {
			return err
		}

		v := value.MapIndex(key)
		err = encoder.writeValue(v)
		if err != nil {
			return err
		}
		encoder.popKey()
	}

	return encoder.setString("")
}

func (encoder *Encoder) writeStruct(v reflect.Value) error {

	err := encoder.writeByte(OBJECT_MARKER)
	if err != nil {
		return err
	}

	if v.CanAddr() {
		index, ok := encoder.ocache[v.Addr().Pointer()]
		if ok {
			index <<= 1
			encoder.writeU29(uint32(index << 1))
			return nil
		}
	}

	err = encoder.writeByte(0x0b)
	if err != nil {
		return err
	}

	err = encoder.setString("")
	if err != nil {
		return err
	}

	t := v.Type()
	switch t.Kind() {
	case reflect.Struct:
		ok := false
		var namer AmfNamer
		if v.CanAddr() {
			namer, ok = v.Addr().Interface().(AmfNamer)
		}
		for i := 0; i < t.NumField(); i++ {

			f := t.Field(i)
			key := f.Name
			if !misc.Exportable(key) {
				continue
			}

			if ok {
				key = namer.GetAmfName(f.Name)
			} else {
				tag := f.Tag.Get("amf")
				if tag != "" {
					key = tag
				}
			}

			if key == "" {
				continue
			}

			encoder.pushKey(key)
			err = encoder.setString(key)
			if err != nil {
				return err
			}

			fv := v.FieldByName(f.Name)
			err = encoder.writeValue(fv)
			if err != nil {
				return err
			}

			encoder.popKey()
		}
	default:
		panic("not a struct")
	}

	return encoder.setString("")
}

func (encoder *Encoder) writeSlice(value reflect.Value) error {

	err := encoder.writeByte(ARRAY_MARKER)
	if err != nil {
		return err
	}

	index, ok := encoder.ocache[value.Pointer()]
	if ok {
		index <<= 1
		encoder.writeU29(uint32(index << 1))
		return nil
	}

	err = encoder.writeU29((uint32(value.Len()) << 1) | 0x01)
	if err != nil {
		return err
	}

	//FIXME 这里未实现ECMA数组

	err = encoder.setString("")
	if err != nil {
		return err
	}

	for i := 0; i < value.Len(); i++ {
		encoder.pushKey(strconv.Itoa(i))
		v := value.Index(i)
		err = encoder.writeValue(v)
		if err != nil {
			return err
		}
		encoder.popKey()
	}

	return nil
}

func (encoder *Encoder) writeValue(v reflect.Value) error {
	switch v.Kind() {
	case reflect.Map:
		return encoder.writeMap(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return encoder.writeUint(v.Uint())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return encoder.writeInt(v.Int())
	case reflect.String:
		return encoder.writeString(v.String())
	case reflect.Array:
		v = v.Slice(0, v.Len())
		return encoder.writeSlice(v)
	case reflect.Slice:
		if v.IsNil() {
			return encoder.writeNull()
		}
		return encoder.writeSlice(v)
	case reflect.Float64, reflect.Float32:
		return encoder.writeFloat(v.Float())
	case reflect.Interface, reflect.Ptr:
		if v.IsNil() {
			return encoder.writeNull()
		}
		v = v.Elem()
		return encoder.writeValue(v)
	case reflect.Struct:
		return encoder.writeStruct(v)
	case reflect.Bool:
		return encoder.writeBool(v.Bool())
	}

	return NewAmfError(encoder.KeyStack, "unsupported type:"+v.Type().String())
}

func (encoder *Encoder) Encode(data interface{}) (err error) {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Ptr {
		return NewAmfError(encoder.KeyStack, "pointer expected")
	}
	return encoder.writeValue(v)
}

func (encoder *Encoder) setString(value string) error {

	index, ok := encoder.scache[value]
	if ok {
		encoder.writeU29(uint32(index << 1))
		return nil
	}

	err := encoder.writeU29(uint32((len(value) << 1) | 0x01))
	if err != nil {
		return err
	}

	if value != "" {
		encoder.scache[value] = len(encoder.scache)
	}
	return encoder.writeBytes([]byte(value))
}

func (encoder *Encoder) writeByte(value byte) error {

	return encoder.writeBytes([]byte{value})
}

func (encoder *Encoder) writeBytes(bytes []byte) error {

	length, err := encoder.writer.Write(bytes)
	if length != len(bytes) || err != nil {
		return NewAmfError(encoder.KeyStack, "write data failed")
	}
	return err
}

func (encoder *Encoder) writeU29(value uint32) error {

	buffer := make([]byte, 0, 4)

	switch {
	case value < 0x80:
		buffer = append(buffer, byte(value))
	case value < 0x4000:
		buffer = append(buffer, byte((value>>7)|0x80))
		buffer = append(buffer, byte(value&0x7f))
	case value < 0x200000:
		buffer = append(buffer, byte((value>>14)|0x80))
		buffer = append(buffer, byte((value>>7)|0x80))
		buffer = append(buffer, byte(value&0x7f))
	case value < 0x20000000:
		buffer = append(buffer, byte((value>>22)|0x80))
		buffer = append(buffer, byte((value>>15)|0x80))
		buffer = append(buffer, byte((value>>7)|0x80))
		buffer = append(buffer, byte(value&0xff))
	default:
		return NewAmfError(encoder.KeyStack, "u29 overflow")
	}

	return encoder.writeBytes(buffer)
}

func NewEncoder(writer io.Writer, byteOrder binary.ByteOrder) *Encoder {

	encoder := new(Encoder)
	encoder.writer = writer
	encoder.byteOrder = byteOrder
	encoder.Reset()
	return encoder
}

/* vim: set ts=4 sw=4 sts=4 tw=100 noet: */
