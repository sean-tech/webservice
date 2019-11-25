package encrypt

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/sean-tech/webservice/fileutils"
	"log"
	"testing"
)

func TestMd5EncryptImpl_EncryptWithTimestamp(t *testing.T) {
	GetMd5().EncryptWithTimestamp([]byte("i am yang"), 0)
}

func TestRsaEncryptImpl_Encrypt(t *testing.T) {

	buf, err := fileutils.ReadFile("/Users/lyra/Desktop/Doc/安全方案/businessS/spubkey.pem")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(buf))

	src := "this is a test word for rsa encrypt"
	encryptData, err := GetRsa().Encrypt(string(buf), []byte(src))
	if err != nil {
		t.Error(err)
	}
	fmt.Println(base64.StdEncoding.EncodeToString(encryptData))
}

func TestRsaEncryptImpl_Decrypt(t *testing.T) {

	buf, err := fileutils.ReadFile("/Users/lyra/Desktop/Doc/安全方案/businessS/sprivkey.pem")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(buf))

	src := "SWD8oHiSHA645ez1isB8fuXy6JLhDgQfDbvWHUUYDswg0qeTV6i3g9dQ/yZBMd0UbjEpmo03D9dS54WMAF4BGVRtkizJiecqxL4Hm6O4hWqSzaQxunIcv2seC5qmJbVLP4SNvv+Y/BQ9k5me9mqS7W0xucb3Jj6U2FqDybU2+9E="
	data, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		t.Error(err)
	}
	decryptData, err := GetRsa().Decrypt(string(buf), data)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(decryptData))
}

func TestRsaEncryptImpl_Sign(t *testing.T) {

	buf, err := fileutils.ReadFile("/Users/lyra/Desktop/Doc/安全方案/businessS/sprivkey.pem")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(buf))

	src := "this is a test word for rsa sign"
	signData, err := GetRsa().Sign(string(buf), []byte(src))
	if err != nil {
		t.Error(err)
	}
	fmt.Println(base64.StdEncoding.EncodeToString(signData))
}

func TestRsaEncryptImpl_Verify(t *testing.T) {

	buf, err := fileutils.ReadFile("/Users/lyra/Desktop/Doc/安全方案/businessS/spubkey.pem")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(buf))

	src := "this is a test word for rsa sign"
	signStr := "KuxIs8BOfmAtHdH0nwE/269WwrDHaX/4sZ4eJUAyTx7+p5CH3w5CVxEbPdLztPKTvIn499Kzg+3q3WHS38oGUzjNVQdE942VxvnEtsOBABlRJlZl4RPIdiLKO6Z8qww1tKd427GHQrYombOLF7jXRmpG9CwD4to27oFWFi0Z1bA="
	signData, err := base64.StdEncoding.DecodeString(signStr)
	if err != nil {
		t.Error(err)
	}
	err = GetRsa().Verify(string(buf), []byte(src), signData)
	if err != nil {
		t.Error(err)
	} else {
		log.Println("verify success")
	}
}

func TestAesEncryptImpl_EncryptCBC(t *testing.T) {

	key := GetAes().GenerateKey()
	fmt.Println(key)
	fmt.Println(hex.EncodeToString(key))
	fmt.Println(base64.StdEncoding.EncodeToString(key))

	origData := []byte("Hello World") // 待加密的数据

	log.Println("------------------ CBC模式 --------------------")
	encrypted, err := GetAes().EncryptCBC(origData, key)
	if err != nil {
		t.Error(err)
	}
	log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	log.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))

	decrypted, err := GetAes().DecryptCBC(encrypted, key)
	if err != nil {
		t.Error(err)
	}
	log.Println("解密结果：", string(decrypted))

	//log.Println("------------------ ECB模式 --------------------")
	//encrypted = AesEncryptECB(origData, key)
	//log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	//log.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	//decrypted = AesDecryptECB(encrypted, key)
	//log.Println("解密结果：", string(decrypted))
	//
	//log.Println("------------------ CFB模式 --------------------")
	//encrypted = AesEncryptCFB(origData, key)
	//log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	//log.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	//decrypted = AesDecryptCFB(encrypted, key)
	//log.Println("解密结果：", string(decrypted))

}