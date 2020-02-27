## s3 api

重新封装[s3](github.com/aws/aws-sdk-go/service/s3)的api接口。

### 安装

> go get -u github.com/zhufuyi/pkg/awss3

<br>

### 使用

使用方式请看[test文件](./awss3_test.go)。

**注：预签名下载功能只支持旧版本的签名协议，现在新创建的bucket暂不支持预签名下载，显示签名不匹配错误，后期需要改进。**
