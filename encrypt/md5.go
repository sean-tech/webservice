package encrypt

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

type md5EncryptImpl struct {}

func (this *md5EncryptImpl) EncryptWithTimestamp(value string) string {
	m := md5.New()
	v := fmt.Sprintf("%s%s", time.Now().String(), value)
	m.Write([]byte(v))
	return hex.EncodeToString(m.Sum(nil))
}

func (this *md5EncryptImpl) Encrypt2WithTimestamp(value string) string {
	m := md5.New()
	m.Write([]byte(value))
	encryptStr := hex.EncodeToString(m.Sum(nil))
	return this.EncryptWithTimestamp(encryptStr)
}