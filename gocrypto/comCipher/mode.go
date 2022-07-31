package comCipher

import "crypto/cipher"

// CipherMode 为不同的工作模式提供了统一的接口来设置填充方式，创建Cipher。
type CipherMode interface {
	SetPadding(padding Padding) CipherMode
	Cipher(block cipher.Block, iv []byte) Cipher
}

type cipherMode struct {
	padding Padding
}

func (c *cipherMode) SetPadding(padding Padding) CipherMode {
	return c
}

func (c *cipherMode) Cipher(block cipher.Block, iv []byte) Cipher {
	return nil
}

type ecbCipherModel cipherMode

func NewECBMode() CipherMode {
	return &ecbCipherModel{padding: NewPKCS57Padding()}
}

func (ecb *ecbCipherModel) SetPadding(padding Padding) CipherMode {
	ecb.padding = padding
	return ecb
}

func (ecb *ecbCipherModel) Cipher(block cipher.Block, iv []byte) Cipher {
	encrypter := NewECBEncrypter(block)
	decrypter := NewECBDecrypter(block)
	return NewBlockCipher(ecb.padding, encrypter, decrypter)
}

type cbcCipherModel cipherMode

func NewCBCMode() CipherMode {
	return &cbcCipherModel{padding: NewPKCS57Padding()}
}

func (cbc *cbcCipherModel) SetPadding(padding Padding) CipherMode {
	cbc.padding = padding
	return cbc
}

func (cbc *cbcCipherModel) Cipher(block cipher.Block, iv []byte) Cipher {
	encrypter := cipher.NewCBCEncrypter(block, iv)
	decrypter := cipher.NewCBCDecrypter(block, iv)
	return NewBlockCipher(cbc.padding, encrypter, decrypter)
}

type cfbCipherModel cipherMode

func NewCFBMode() CipherMode {
	return &ofbCipherModel{}
}

func (cfb *cfbCipherModel) Cipher(block cipher.Block, iv []byte) Cipher {
	encrypter := cipher.NewCFBEncrypter(block, iv)
	decrypter := cipher.NewCFBDecrypter(block, iv)
	return NewStreamCipher(encrypter, decrypter)
}

type ofbCipherModel struct {
	cipherMode
}

func NewOFBMode() CipherMode {
	return &ofbCipherModel{}
}

func (ofb *ofbCipherModel) Cipher(block cipher.Block, iv []byte) Cipher {
	encrypter := cipher.NewOFB(block, iv)
	decrypter := cipher.NewOFB(block, iv)
	return NewStreamCipher(encrypter, decrypter)
}

type ctrCipherModel struct {
	cipherMode
}

func NewCTRMode() CipherMode {
	return &ctrCipherModel{}
}

func (ctr *ctrCipherModel) Cipher(block cipher.Block, iv []byte) Cipher {
	encrypter := cipher.NewCTR(block, iv)
	decrypter := cipher.NewCTR(block, iv)
	return NewStreamCipher(encrypter, decrypter)
}
