## s3 api

重新封装[s3](github.com/aws/aws-sdk-go/service/s3)的api接口。

<br>

## 安装

> go get -u github.com/zhufuyi/pkg/awss3

<br>

## 使用示例

```go
    // 初始化
    s3Cli, err = awss3.NewAwsS3(bucket, basePath, region, credentialsFile)

    // 上传本地文件
    url, err := s3Cli.UploadFromFile(localFile)
    // 上传内容
    url, err := s3Cli.UploadFromReader(f, localFile)

    // 下载文件
    n, err := s3Cli.DownloadToFile(awsFile, localFile)

    // 删除文件
    err := s3Cli.DeleteFile(awsFile)

    // 判断文件是否存在
    err := s3Cli.CheckFileIsExist(awsFile)

    // 获取访问s3资源的签名url，给第三方上传和下载资源使用
    url, err := s3Cli.GetPreSignedURL(DownloadMethod, awsFile, 1000)
```

