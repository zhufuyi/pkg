package awss3

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	// StorageClass 存储类型
	StorageClass = "STANDARD_IA"

	// AccessKeyIDFieldName aws配置文件credentials的访问id字段名称
	AccessKeyIDFieldName = "aws_access_key_id"
	// SecretAccessKeyFieldName aws配置文件credentials的访问密钥字段名称
	SecretAccessKeyFieldName = "aws_secret_access_key"

	// UploadMethod 预签名上传
	UploadMethod = "put"
	// DownloadMethod 预签名下载
	DownloadMethod = "get"
)

// AwsS3 封装aws s3对象
type AwsS3 struct {
	Bucket          string // aws的存储桶
	BasePath        string // 结尾必须以路径分隔符号/结尾
	Region          string // 区域
	accessKeyID     string // aws访问id
	secretAccessKey string // aws访问密钥

	Session *session.Session // 会话
	S3      *s3.S3           // s3对象
}

// NewAwsS3 使用初始化s3，建议优先使用.aws/credentials文件，如果不配置文件，可以传递aws_access_key_id和aws_secret_access_key值
func NewAwsS3(bucket string, basePath string, region string, credentialsFile string, credentialsValues ...string) (*AwsS3, error) {
	var accessKeyID, secretAccessKey string
	var err error

	// 如果没有配置.aws/credentials文件，尝试从可变参数中获取值，第一个参数为aws_access_key_id的值，第二个参数为aws_secret_access_key的值
	accessKeyID, secretAccessKey, _ = getCredentialsValues(credentialsFile)
	if accessKeyID == "" || secretAccessKey == "" {
		if len(credentialsValues) == 2 {
			accessKeyID = credentialsValues[0]
			secretAccessKey = credentialsValues[1]
		}
	}

	switch "" {
	case region:
		return nil, errors.New("region value is empty")
	case bucket:
		return nil, errors.New("bucket value is empty")
	case accessKeyID:
		return nil, errors.New("accessKeyID value is empty")
	case secretAccessKey:
		return nil, errors.New("secretAccessKey value is empty")
	}

	creds := credentials.NewStaticCredentials(accessKeyID, secretAccessKey, "")
	cfg := aws.NewConfig().WithRegion(region).WithCredentials(creds)
	sess, err := session.NewSession(cfg)
	if err != nil {
		return nil, err
	}

	// 确保basePath为path/或空两种方式
	if basePath != "" {
		if basePath[0] == '/' {
			basePath = strings.TrimLeft(basePath, "/")
		}
		if basePath != "" && basePath[len(basePath)-1] != '/' {
			basePath += "/"
		}
	}

	awss3 := &AwsS3{
		Bucket:          bucket,
		BasePath:        basePath,
		Region:          region,
		accessKeyID:     accessKeyID,
		secretAccessKey: secretAccessKey,
		Session:         sess,
		S3:              s3.New(sess),
	}

	return awss3, nil
}

// UploadFromFile 把文件上传到s3，s3上的文件名和上传的文件名一致
func (a *AwsS3) UploadFromFile(localFile string) (string, error) {
	f, err := os.Open(localFile)
	if err != nil {
		return "", err
	}
	defer f.Close() //nolint

	uploader := s3manager.NewUploader(a.Session)
	uploadInput := &s3manager.UploadInput{
		Bucket:       aws.String(a.Bucket),
		Key:          aws.String(a.BasePath + filepath.Base(localFile)),
		StorageClass: aws.String(StorageClass),
		Body:         f,
	}

	output, err := uploader.Upload(uploadInput)
	if err != nil {
		return "", err
	}

	return output.Location, err
}

// UploadFromReader 读取reader内容上传到s3，可以指定以文件名
func (a *AwsS3) UploadFromReader(reader io.Reader, awsFile string) (string, error) {
	uploader := s3manager.NewUploader(a.Session)
	uploadInput := &s3manager.UploadInput{
		Bucket:       aws.String(a.Bucket),
		Key:          aws.String(a.BasePath + filepath.Base(awsFile)),
		StorageClass: aws.String(StorageClass),
		Body:         reader,
	}

	output, err := uploader.Upload(uploadInput)
	if err != nil {
		return "", err
	}

	return output.Location, err
}

// DownloadToFile 从S3下载文件到本地
func (a *AwsS3) DownloadToFile(awsFile string, localFile string) (int64, error) {
	f, err := os.Create(localFile)
	if err != nil {
		return 0, err
	}
	defer f.Close() //nolint

	downloader := s3manager.NewDownloader(a.Session)
	input := &s3.GetObjectInput{
		Bucket: aws.String(a.Bucket),
		Key:    aws.String(a.BasePath + filepath.Base(awsFile)),
	}

	return downloader.Download(f, input)
}

// CheckFileIsExist 检查对象是否存在
func (a *AwsS3) CheckFileIsExist(awsFile string) error {
	input := &s3.GetObjectInput{
		Bucket: aws.String(a.Bucket),
		Key:    aws.String(a.BasePath + awsFile),
	}
	_, err := a.S3.GetObject(input)

	return err
}

// DeleteFile 从s3删除文件
func (a *AwsS3) DeleteFile(awsFile string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(a.Bucket),
		Key:    aws.String(a.BasePath + awsFile),
	}
	_, err := a.S3.DeleteObject(input)

	return err
}

// 生成可以上传或下载的s3预签名url
func (a *AwsS3) genPreSignedURL(method string, uri string, expiry time.Duration) (*string, time.Time, error) {
	bucket, key, err := parseS3URI(uri)
	if err != nil {
		return nil, time.Now().UTC(), err
	}

	var req *request.Request
	if method == "put" {
		req, _ = a.S3.PutObjectRequest(&s3.PutObjectInput{
			Bucket: bucket,
			Key:    key,
		})
	} else if method == "get" {
		req, _ = a.S3.GetObjectRequest(&s3.GetObjectInput{
			Bucket: bucket,
			Key:    key,
		})
	}

	if req == nil {
		return nil, time.Now().UTC(), errors.New("Unable to create a request")
	}

	s3url, err := req.Presign(expiry)
	if err != nil {
		return nil, time.Now().UTC(), err
	}

	return &s3url, time.Now().UTC().Add(expiry), nil
}

// GetPreSignedURL 获取访问s3资源的签名url，给第三方上传和下载资源使用
func (a *AwsS3) GetPreSignedURL(method string, awsFile string, expirySeconds int) (string, error) {
	if method != DownloadMethod && method != UploadMethod {
		return "", errors.New("method does not exist, must use 'get' or 'put'")
	}
	urlStr := "s3://" + a.Bucket + "/" + a.BasePath + awsFile
	// 默认1800秒
	if expirySeconds <= 10 {
		expirySeconds = 1800
	}
	expiry := time.Second * time.Duration(expirySeconds)

	s3Url, _, err := a.genPreSignedURL(method, urlStr, expiry)
	if err != nil {
		return "", err
	}

	return *s3Url, nil
}

// WithBucketAndBasePath 切换新的S3存储桶和基础路径，复制新的s3对象，但不会影响原s3，如果不填写bucket则使用默认值
func (a *AwsS3) WithBucketAndBasePath(currentBasePath string, currentBucket ...string) *AwsS3 {
	if currentBasePath != "" {
		if currentBasePath[0] == '/' {
			currentBasePath = strings.TrimLeft(currentBasePath, "/")
		}
		if currentBasePath != "" && currentBasePath[len(currentBasePath)-1] != '/' {
			currentBasePath += "/"
		}
	}

	currentS3 := *a
	if len(currentBucket) == 1 && currentBucket[0] != "" {
		currentS3.Bucket = currentBucket[0]
	}
	currentS3.BasePath = currentBasePath

	return &currentS3
}

// -------------------------------------------------------------------------------------------------

func preProcess(line string) []string {
	line = strings.Replace(line, "\n", "", -1)
	line = strings.Replace(line, "\r", "", -1)
	line = strings.Replace(line, " ", "", -1)
	return strings.Split(line, "=")
}

func getCredentialsValues(credentialsFile string) (string, string, error) {
	var accessKeyID, secretAccessKey string

	if credentialsFile == "" {
		errMsg := "credentialsFile value is empty, use specified parameters"
		fmt.Println(errMsg)
		return "", "", errors.New(errMsg)
	}

	f, err := os.Open(credentialsFile)
	if err != nil {
		return "", "", err
	}
	defer f.Close() //nolint

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
		if err != nil || io.EOF == err {
			break
		}

		if strings.Contains(line, AccessKeyIDFieldName) {
			ss := preProcess(line)
			if len(ss) == 2 && ss[0] == AccessKeyIDFieldName && ss[1] != "" {
				accessKeyID = ss[1]
			}
		} else if strings.Contains(line, SecretAccessKeyFieldName) {
			ss := preProcess(line)
			if len(ss) == 2 && ss[0] == SecretAccessKeyFieldName && ss[1] != "" {
				secretAccessKey = ss[1]
			}
		}
		if accessKeyID != "" && secretAccessKey != "" { // 只读取第一对参数值，如果有其他aws帐户下值会被忽略
			break
		}
	}

	return accessKeyID, secretAccessKey, nil
}

func parseS3URI(uri string) (*string, *string, error) {
	parsed, err := url.Parse(uri)
	if err != nil {
		return nil, nil, errors.New("Unable to parse S3 URL: " + uri)
	}

	bucket := parsed.Host
	key := parsed.Path[1:]
	return &bucket, &key, nil
}

// -------------------------------------------------------------------------------------------------

var myS3 *AwsS3

// InitS3 初始化s3
func InitS3(bucket string, basePath string, region string, credentialsFile string, credentialsValues ...string) error {
	var err error
	myS3, err = NewAwsS3(bucket, basePath, region, credentialsFile, credentialsValues...)
	return err
}

// GetS3 获取s3对象，如果未初始化就使用会发生panic
func GetS3() *AwsS3 {
	if myS3 == nil {
		panic("s3 is nil, not initialized yet.")
	}

	return myS3
}
