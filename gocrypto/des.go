package gocrypto

// 对称加密DES，是目前最为流行的加密算法之一，逐渐被AES替换。

import (
	"encoding/hex"

	"github.com/zhufuyi/pkg/gocrypto/comCipher"
)

// DesEncrypt des加密byte，返回的密文未经过转码
func DesEncrypt(rawData []byte, opts ...DesOption) ([]byte, error) {
	o := defaultDesOptions()
	o.apply(opts...)

	return desEncrypt(o.mode, rawData, o.desKey)
}

// DesDecrypt des解密byte，参数输入未经过转码的密文
func DesDecrypt(cipherData []byte, opts ...DesOption) ([]byte, error) {
	o := defaultDesOptions()
	o.apply(opts...)

	return desDecrypt(o.mode, cipherData, o.desKey)
}

// DesEncryptHex des加密string，返回的密文已经转码
func DesEncryptHex(rawData string, opts ...DesOption) (string, error) {
	o := defaultDesOptions()
	o.apply(opts...)

	cipherData, err := desEncrypt(o.mode, []byte(rawData), o.desKey)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(cipherData), nil
}

// DesDecryptHex des解密string，参数输入已经转码的密文字符串
func DesDecryptHex(cipherStr string, opts ...DesOption) (string, error) {
	o := defaultDesOptions()
	o.apply(opts...)

	cipherData, err := hex.DecodeString(cipherStr)
	if err != nil {
		return "", err
	}

	rawData, err := desDecrypt(o.mode, cipherData, o.desKey)
	if err != nil {
		return "", err
	}

	return string(rawData), nil
}

func desEncrypt(mode string, rawData []byte, key []byte) ([]byte, error) {
	cipherMode, err := getCipherMode(mode)
	if err != nil {
		return nil, err
	}

	cip, err := comCipher.NewDESWith(key, cipherMode)
	if err != nil {
		return nil, err
	}

	return cip.Encrypt(rawData), nil
}

func desDecrypt(mode string, cipherData []byte, key []byte) ([]byte, error) {
	cipherMode, err := getCipherMode(mode)
	if err != nil {
		return nil, err
	}

	cip, err := comCipher.NewDESWith(key, cipherMode)
	if err != nil {
		return nil, err
	}

	return cip.Decrypt(cipherData), nil
}
