package encrypt

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"time"
)

type md5EncryptImpl struct {}

func (this *md5EncryptImpl) Encrypt(value []byte) string {
	m := md5.New()
	m.Write(value)
	return hex.EncodeToString(m.Sum(nil))
}

func (this *md5EncryptImpl) EncryptWithTimestamp(value []byte, timestamp int64) string {
	if timestamp == 0 {
		timestamp = time.Now().Unix()
	}
	var buf bytes.Buffer
	buf.Write(value)
	buf.WriteString(strconv.FormatInt(timestamp, 10))
	m := md5.New()
	m.Write(buf.Bytes())
	return hex.EncodeToString(m.Sum(nil))
}

