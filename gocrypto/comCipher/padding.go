package comCipher

import "bytes"

// Padding 为各种填充方式提供了统一的接口来填充/还原数据。
type Padding interface {
	Padding(src []byte, blockSize int) []byte
	UnPadding(src []byte) []byte
}

type padding struct{}

type pkcs57Padding padding

func NewPKCS57Padding() Padding {
	return &pkcs57Padding{}
}

func (p *pkcs57Padding) Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func (p *pkcs57Padding) UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}
