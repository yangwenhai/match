// amf库实现了对于amf3的编解码功能
package amf

import (
	"encoding/binary"
	"io"
)

func Encode(writer io.Writer, byteOrder binary.ByteOrder, data interface{}) error {
	encoder := NewEncoder(writer, byteOrder)
	return encoder.Encode(data)
}

func Decode(reader io.Reader, byteOrder binary.ByteOrder, data interface{}) error {
	decoder := NewDecoder(reader, byteOrder)
	return decoder.Decode(data)
}

/* vim: set ts=4 sw=4 sts=4 tw=100 noet: */
