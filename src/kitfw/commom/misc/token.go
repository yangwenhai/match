package misc

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"
)

func NextTokenAsNumber() uint64 {
	now := time.Now().UnixNano()
	suffix := rand.Int63()
	str := fmt.Sprintf("%d#%d", now, suffix)

	hash := md5.New()
	ret := hash.Sum([]byte(str))

	buffer := bytes.NewBuffer(ret)
	var high, low uint64
	binary.Read(buffer, binary.BigEndian, &high)
	binary.Read(buffer, binary.BigEndian, &low)

	token := ((high ^ low) % 1e8) * 100
	if token < 1e9 {
		token += 1e9
	}

	return token
}

func NextTokenAsString() string {
	return fmt.Sprintf("%d", NextTokenAsNumber())
}

/* vim: set ts=4 sw=4 sts=4 tw=100 noet: */
