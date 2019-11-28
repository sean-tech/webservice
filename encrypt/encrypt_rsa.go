package encrypt

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

type rsaEncryptImpl struct {}

/**
 * 公钥加密
 */
func (this *rsaEncryptImpl) Encrypt(publicKey string, data []byte) ([]byte, error) {
	publicKeyBytes := []byte(publicKey)
	block, _ := pem.Decode(publicKeyBytes)
	if block == nil {
		return nil,errors.New("public key error")
	}
	pukI, err := x509.ParsePKIXPublicKey([]byte(block.Bytes))
	if err != nil {
		return nil,err
	}
	pubKey := pukI.(*rsa.PublicKey)
	//避免数据过长报错，故分段加密
	partLen := pubKey.N.BitLen() / 8 - 11
	chunks := split(data, partLen)
	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		bytes, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, chunk)
		if err != nil {
			return nil, err
		}
		buffer.Write(bytes)
	}
	return buffer.Bytes(), nil
}

//用进行操作
/**
 * 私钥解密(PKCS8)
 */
func (this *rsaEncryptImpl) Decrypt(privateKey string, data []byte) ([]byte, error) {
	pri_Key_bye := []byte(privateKey)

	blockPri, _ := pem.Decode(pri_Key_bye)
	if blockPri == nil {
		return nil, errors.New("private key error")
	}

	prkI, err := x509.ParsePKCS8PrivateKey([]byte(blockPri.Bytes))
	//priKey, err := x509.ParsePKCS1PrivateKey([]byte(blockPri.Bytes))
	if err != nil {
		return nil,err
	}
	priKey := prkI.(*rsa.PrivateKey)

	partLen := priKey.N.BitLen() / 8
	chunks := split(data, partLen)
	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, priKey, chunk)
		if err != nil {
			return nil, err
		}
		buffer.Write(decrypted)
	}

	return buffer.Bytes(), nil
}


//用进行操作
/**
 * 私钥签名(PKCS8)
 */
func (this *rsaEncryptImpl) Sign(privateKey string, data []byte) ([]byte, error) {
	pri_Key_bye := []byte(privateKey)

	blockPri, _ := pem.Decode(pri_Key_bye)
	if blockPri == nil {
		return nil, errors.New("private key error")
	}

	prkI, err := x509.ParsePKCS8PrivateKey([]byte(blockPri.Bytes))
	if err != nil {
		return nil,err
	}
	priKey := prkI.(*rsa.PrivateKey)

	//hashedStr := GetMd5().Encrypt(data)
	//hashed,_ := hex.DecodeString(hashedStr)
	h := sha1.New()
	h.Write(data)
	hashed := h.Sum(nil)

	return rsa.SignPKCS1v15(rand.Reader, priKey, crypto.SHA1, hashed)
}

/**
 * 公钥验签
 */
func (this *rsaEncryptImpl) Verify(publicKey string, data []byte, signedData []byte) error {
	pub_Key_bye := []byte(publicKey)

	blockPri, _ := pem.Decode(pub_Key_bye)
	if blockPri == nil {
		return errors.New("public key error")
	}

	pukI, err := x509.ParsePKIXPublicKey([]byte(blockPri.Bytes))
	if err != nil {
		return err
	}
	pubKey := pukI.(*rsa.PublicKey)

	h := sha1.New()
	h.Write(data)
	hashed := h.Sum(nil)

	//hashedStr := GetMd5().Encrypt(data)
	//hashed,_ := hex.DecodeString(hashedStr)

	return rsa.VerifyPKCS1v15(pubKey, crypto.SHA1, hashed, signedData)

}

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}