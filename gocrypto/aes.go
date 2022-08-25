package gocrypto

// 对称加密AES，高级加密标准，安全度最高，已逐渐替代DES成为新一代对称加密的标准

import (
	"encoding/hex"
	"errors"

	"github.com/zhufuyi/pkg/gocrypto/comCipher"
)

// AesEncrypt aes加密byte，返回的密文未经过转码
func AesEncrypt(rawData []byte, opts ...AesOption) ([]byte, error) {
	o := defaultAesOptions()
	o.apply(opts...)

	return aesEncrypt(o.mode, rawData, o.aesKey)
}

// AesDecrypt aes解密byte，参数输入未经过转码的密文
func AesDecrypt(cipherData []byte, opts ...AesOption) ([]byte, error) {
	o := defaultAesOptions()
	o.apply(opts...)

	return aesDecrypt(o.mode, cipherData, o.aesKey)
}

// AesEncryptHex aes加密string，返回的密文已经转码
func AesEncryptHex(rawData string, opts ...AesOption) (string, error) {
	o := defaultAesOptions()
	o.apply(opts...)

	cipherData, err := aesEncrypt(o.mode, []byte(rawData), o.aesKey)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(cipherData), nil
}

// AesDecryptHex aes解密string，参数输入已经转码的密文字符串
func AesDecryptHex(cipherStr string, opts ...AesOption) (string, error) {
	o := defaultAesOptions()
	o.apply(opts...)

	cipherData, err := hex.DecodeString(cipherStr)
	if err != nil {
		return "", err
	}

	rawData, err := aesDecrypt(o.mode, cipherData, o.aesKey)
	if err != nil {
		return "", err
	}

	return string(rawData), nil
}

func getCipherMode(mode string) (comCipher.CipherMode, error) {
	var cipherMode comCipher.CipherMode
	switch mode {
	case modeECB:
		cipherMode = comCipher.NewECBMode()
	case modeCBC:
		cipherMode = comCipher.NewCBCMode()
	case modeCFB:
		cipherMode = comCipher.NewCFBMode()
	case modeCTR:
		cipherMode = comCipher.NewCTRMode()
	default:
		return nil, errors.New("unknown mode = " + mode)
	}

	return cipherMode, nil
}

func aesEncrypt(mode string, rawData []byte, key []byte) ([]byte, error) {
	cipherMode, err := getCipherMode(mode)
	if err != nil {
		return nil, err
	}

	cip, err := comCipher.NewAESWith(key, cipherMode)
	if err != nil {
		return nil, err
	}

	return cip.Encrypt(rawData), nil
}

func aesDecrypt(mode string, cipherData []byte, key []byte) ([]byte, error) {
	cipherMode, err := getCipherMode(mode)
	if err != nil {
		return nil, err
	}

	cip, err := comCipher.NewAESWith(key, cipherMode)
	if err != nil {
		return nil, err
	}

	return cip.Decrypt(cipherData), nil
}
