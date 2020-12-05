package encrypt

import "sync"

type IMd5Encrypt interface {
	Encrypt(value []byte) string
	EncryptWithTimestamp(value []byte, timestamp int64) string
}

type IRsaEncrypt interface {
	Encrypt(publicKey string, data []byte) ([]byte, error)
	Decrypt(privateKey string, data []byte) ([]byte, error)
	Sign(privateKey string, data []byte) ([]byte, error)
	Verify(publicKey string, data []byte, signedData []byte) error
}

type IAesEncrypt interface {
	EncryptCBC(origData []byte, key []byte) ([]byte, error)
	DecryptCBC(encrypted []byte, key []byte) ([]byte, error)
	GenerateKey() []byte
}

var (
	_md5Once 		sync.Once
	_md5Instance 	IMd5Encrypt

	_rsaOnce 		sync.Once
	_rsaInstance 	IRsaEncrypt

	_aesOnce 		sync.Once
	_aesInstance 	IAesEncrypt
)

func GetMd5() IMd5Encrypt {
	_md5Once.Do(func() {
		_md5Instance = new(md5EncryptImpl)
	})
	return _md5Instance
}

func GetRsa() IRsaEncrypt {
	_rsaOnce.Do(func() {
		_rsaInstance = new(rsaEncryptImpl)
	})
	return _rsaInstance
}

func GetAes() IAesEncrypt {
	_aesOnce.Do(func() {
		_aesInstance = new(aesEncryptImpl)
	})
	return _aesInstance
}