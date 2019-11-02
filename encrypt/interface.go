package encrypt

import "sync"

type IMd5Encrypt interface {
	EncryptWithTimestamp(value string) string
	Encrypt2WithTimestamp(value string) string
}

var (
	md5Once 		sync.Once
	md5Instance 	IMd5Encrypt
)

func GetMd5Instance() IMd5Encrypt {
	md5Once.Do(func() {
		md5Instance = new(md5EncryptImpl)
	})
	return md5Instance
}