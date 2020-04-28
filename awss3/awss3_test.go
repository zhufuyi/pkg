package awss3

import (
	"fmt"
	"os"
	"testing"
)

var (
	testS3 *AwsS3
	err    error

	bucket   = "mybucket"
	basePath = "/test/"

	region          = "ap-northeast-1"
	credentialsFile = "./credentials"
	accessKeyID     = "xxxxxx"
	secretAccessKey = "xxxxxx"
)

func init() {
	testS3, err = NewAwsS3(bucket, basePath, region, credentialsFile)
	if err != nil {
		panic(err)
	}
}

func TestNewAwsS3(t *testing.T) {
	// 使用配置文件初始化
	as3, err := NewAwsS3(bucket, basePath, region, credentialsFile)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(*as3)

	// 使用参数初始化
	//as3, err = NewAwsS3(bucket, basePath, region, "", accessKeyID, secretAccessKey)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//pp.Println(as3)
}

func TestAwsS3_UploadFromFile(t *testing.T) {
	//localFile := "./uploadTest1.txt"
	//localFile := "uploadTest2.jpg"
	localFile := "uploadTest3.csv"
	url, err := testS3.UploadFromFile(localFile)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("upload success, url =", url)
}

func TestAwsS3_UploadFromReader(t *testing.T) {
	localFile := "./uploadTest4.zip"
	f, err := os.Open(localFile)
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()

	url, err := testS3.UploadFromReader(f, localFile)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("upload success, url =", url)
}

func TestAwsS3_DownloadToFile(t *testing.T) {
	awsFile := "uploadTest1.txt"
	localFile := "./download/" + awsFile
	n, err := testS3.DownloadToFile(awsFile, localFile)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println("download file success, size =", n)
}

func TestAwsS3_CheckFileIsExist(t *testing.T) {
	awsFile := "uploadTest1.txt"
	err := testS3.CheckFileIsExist(awsFile)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(awsFile, "is exist")
}

func TestAwsS3_DeleteFile(t *testing.T) {
	awsFile := "uploadTest1.txt"
	err := testS3.DeleteFile(awsFile)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("delete file %s success.\n", awsFile)
}

func TestAwsS3_GetPreSignedURL(t *testing.T) {
	errMsg := ""
	awsFiles := []string{"uploadTest1.txt", "uploadTest2.jpg", "uploadTest3.csv", "uploadTest4.zip"}
	for _, awsFile := range awsFiles {
		url, err := testS3.GetPreSignedURL(DownloadMethod, awsFile, 1000)
		if err != nil {
			errMsg += err.Error() + "\n"
			continue
		}
		fmt.Println(url)
	}

	if errMsg != "" {
		t.Error(errMsg)
	}
}

func TestAwsS3_WithBucketPath(t *testing.T) {
	currentS3 := testS3.WithBucketAndBasePath("/myPath/")
	fmt.Println(testS3.Bucket, testS3.BasePath)
	fmt.Println(currentS3.Bucket, currentS3.BasePath)
	fmt.Println(testS3.S3, testS3.Session)
	fmt.Println(currentS3.S3, currentS3.Session)

	awsFile := "uploadTest1.txt"
	err := testS3.CheckFileIsExist(awsFile)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(awsFile, "is exist")
	}

	awsFile = "cluster-log-1.gz"
	err = currentS3.WithBucketAndBasePath("/2019/04/05/", "dev-k8s-log").CheckFileIsExist(awsFile)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(awsFile, "is exist")
	}
}
