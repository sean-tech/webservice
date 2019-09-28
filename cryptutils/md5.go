package cryptutils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

type md5Service string
var Md5 *md5Service = new(md5Service)

func (this *md5Service) EncodeWithTimestamp(value string) string {
	m := md5.New()
	v := fmt.Sprintf("%s%s", time.Now().String(), value)
	m.Write([]byte(v))
	return hex.EncodeToString(m.Sum(nil))
}

func (this *md5Service) Encode2WithTimestamp(value string) string {
	m := md5.New()
	m.Write([]byte(value))
	encryptStr := hex.EncodeToString(m.Sum(nil))
	return Md5.EncodeWithTimestamp(encryptStr)
}