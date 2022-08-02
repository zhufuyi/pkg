## gocrypto

封装的常用加密和解密包，包括hash、aes、des、rsa。

<br>

## 安装

> go get -u github.com/zhufuyi/pkg/gocrypto

<br>

## 使用示例

### hash 单向加密

共有几个常用Md5、Sha1、Sha256、Sha512函数和Hash集合，使用示例：

```go
    var hashRawData = []byte("hash|abcdefghijklmnopqrstuvwxyz1234567890")

    // 常用的几个独立的hash函数
    gocrypto.Md5(hashRawData)
    gocrypto.Sha1(hashRawData)
    gocrypto.Sha256(hashRawData)
    gocrypto.Sha512(hashRawData)

    // hash集合，根据哈希类型指定执行对应哈希函数
    gocrypto.Hash(crypto.MD5, hashRawData)
    gocrypto.Hash(crypto.SHA3_224, hashRawData)
    gocrypto.Hash(crypto.SHA256, hashRawData)
    gocrypto.Hash(crypto.SHA3_224, hashRawData)
    gocrypto.Hash(crypto.BLAKE2s_256, hashRawData)
    // ..... 支持crypto包hash类型
```

<br>

### AES 加密解密

AES(`Advanced Encryption Standard`)高级加密标准，旨在取代`DES`，共有四种分组加密模式：ECB CBC CFB CTR。

共有四个函数AesEncrypt、AesDecrypt、AesEncryptHex、AesDecryptHex，使用示例：

```go
    var (
        aesRawData = []byte("aes|abcdefghijklmnopqrstuvwxyz1234567890")
        aesKey     = []byte("aesKey1234567890aesKey1234567890")
    )

    // (1) AesEncrypt和AesDecrypt函数的参数有默认值：模式=ECB，可以修改为CBC CFB CTR，默认值key，可以自定义修改

    // 默认模式ECB，默认key
    cypherData, _ := gocrypto.AesEncrypt(aesRawData) // 加密，返回密文未经过转码
    raw, _ := gocrypto.AesDecrypt(cypherData) // 解密，返回原文

    // 默认模式ECB，自定义key，key长度必须是16、24、32其中一个
    cypherData, _ := gocrypto.AesEncrypt(aesRawData, WithAesKey(aesKey)) // 加密，返回密文未经过转码
    raw, _ := gocrypto.AesDecrypt(cypherData, WithAesKey(aesKey)) // 解密，返回原文

    // 模式CTR，默认key
    cypherData, _ := gocrypto.AesEncrypt(aesRawData, WithAesModeCTR()) // 加密，返回密文未经过转码
    raw, _ := gocrypto.AesDecrypt(cypherData, WithAesModeCTR()) // 解密，返回原文

    // 模式CBC，自定义key，key长度必须是16、24、32其中一个
    cypherData, _ := gocrypto.AesEncrypt(aesRawData, WithAesModeECB(), WithAesKey(aesKey)) // 加密，返回密文未经过转码
    raw, _ := gocrypto.AesDecrypt(cypherData, WithAesModeECB(), WithAesKey(aesKey))        // 解密，返回原文


    // (2) AesEncryptHex和AesDecryptHex函数，这两个函数的密文是经过hex转码，使用方式与AesEncrypt、AesDecrypt完全一样。

```
<br>

### DES 加密解密

DES(`Data Encryption Standard`)数据加密标准，是目前最为流行的加密算法之一  ，共有四种分组加密模式：ECB CBC CFB CTR。

共有四个函数DesEncrypt、DesDecrypt、DesEncryptHex、DesDecryptHex，使用示例：

```go
    var (
        desRawData = []byte("des|abcdefghijklmnopqrstuvwxyz1234567890")
        desKey     = []byte("desKey1234567890desKey1234567890")
    )

    // (1) DesEncrypt和DesDecrypt函数的参数有默认值：模式=ECB，可以修改为CBC CFB CTR，默认值key，可以自定义修改

    // 默认模式ECB，默认key
    cypherData, _ := gocrypto.DesEncrypt(desRawData) // 加密，返回密文未经过转码
    raw, _ := gocrypto.DesDecrypt(cypherData) // 解密，返回原文

    // 默认模式ECB，自定义key，key长度必须是16、24、32其中一个
    cypherData, _ := gocrypto.DesEncrypt(desRawData, WithDesKey(desKey)) // 加密，返回密文未经过转码
    raw, _ := gocrypto.DesDecrypt(cypherData, WithDesKey(desKey)) // 解密，返回原文

    // 模式CTR，默认key
    cypherData, _ := gocrypto.DesEncrypt(desRawData, WithDesModeCTR()) // 加密，返回密文未经过转码
    raw, _ := gocrypto.DesDecrypt(cypherData, WithDesModeCTR()) // 解密，返回原文

    // 模式CBC，自定义key，key长度必须是16、24、32其中一个
    cypherData, _ := gocrypto.DesEncrypt(desRawData, WithDesModeECB(), WithDesKey(desKey)) // 加密，返回密文未经过转码
    raw, _ := gocrypto.DesDecrypt(cypherData, WithDesModeECB(), WithDesKey(desKey))        // 解密，返回原文


    // (2) DesEncryptHex和DesDecryptHex函数，这两个函数的密文是经过hex转码，使用方式与DesEncrypt、DesDecrypt完全一样。

```

<br>

### RSA非对称加密

#### RSA加密和解密

公钥用来加密，私钥用来解密，例如别人用公钥加密加密信息发送给你，你有私钥可以解密出信息内容。加密解密共有4个函数RsaEncrypt、RsaDecrypt、RsaEncryptHex、RsaDecryptHex，使用示例：

```go
    var (
        rsaRawData = []byte("rsa|abcdefghijklmnopqrstuvwxyz1234567890")
    )

    // (1) RsaEncrypt和RsaDecrypt函数的参数有默认值：密钥对格式=PKCS#1，可以修改为PKCS#8

    // 默认密钥对PKCS#1
    cypherData, _ := gocrypto.RsaEncrypt(publicKey, rsaRawData) // 加密，返回密文未经过转码
    raw, _ := gocrypto.RsaDecrypt(privateKey, cypherData) // 解密，返回原文

    // 密钥对PKCS#8
    cypherData, _ := gocrypto.RsaEncrypt(publicKey, rsaRawData, WithRsaFormatPKCS8()) // 加密，返回密文未经过转码
    raw, _ := gocrypto.RsaDecrypt(privateKey, cypherData, WithRsaFormatPKCS8()) // 解密，返回原文


    // (2) RsaEncryptHex和RsaDecryptHex，这两个函数的密文是经过hex转码的，使用方式与RsaEncrypt、RsaDecrypt完全一样。
```

<br>

#### RSA签名和验签

私钥用来签名，公钥用来验签，例如你用私钥对身份签名，别人通过公钥对签名验证得到你身份是否可信任的。签名和验签共有四个函数RsaSign、RsaVerify、RsaSignBase64、RsaVerifyBase64，使用示例：

```go
    var (
        rsaRawData = []byte("rsa|abcdefghijklmnopqrstuvwxyz1234567890")
    )

    // (1) RsaEncrypt和RsaDecrypt函数的参数有默认值：密钥对格式=PKCS#1，可以修改为PKCS#8，默认值哈希=sha1，可以修改为多种哈希类型

    // 默认密钥对PKCS#1，默认哈希sha1
    signData, _ := gocrypto.RsaEncrypt(privateKey, rsaRawData) // 签名，返回密文未经过转码
    err := gocrypto.RsaDecrypt(publicKey, rsaRawData, signData) // 验签

    // 默认密钥对PKCS#1，使用哈希sha256
    signData, _ := gocrypto.RsaEncrypt(privateKey, rsaRawData, WithRsaHashTypeSha256()) // 签名，返回密文未经过转码
    err := gocrypto.RsaDecrypt(publicKey, rsaRawData, signData, WithRsaHashTypeSha256()) // 验签

    // 密钥对PKCS#8，默认哈希sha1
    signData, _ := gocrypto.RsaEncrypt(privateKey, rsaRawData, WithRsaFormatPKCS8()) // 签名，返回密文未经过转码
    err := gocrypto.RsaDecrypt(publicKey, rsaRawData, signData, WithRsaFormatPKCS8()) // 验签

    // 密钥对PKCS#8，哈希sha512
    signData, _ := gocrypto.RsaEncrypt(privateKey, rsaRawData, WithRsaFormatPKCS8(), WithRsaHashTypeSha512()) // 签名，返回密文未经过转码
    err := gocrypto.RsaDecrypt(publicKey, rsaRawData, signData, WithRsaFormatPKCS8(), WithRsaHashTypeSha512()) // 验签


    // (2) RsaSignBase64和RsaVerifyBase64这两个函数的密文是经过base64转码的，使用方式与RsaEncrypt、RsaDecrypt完全一样。

```
