package awss3

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/zhufuyi/pkg/utils"
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
		t.Log(err)
		return
	}
	fmt.Println(*as3)

	// 使用参数初始化
	as3, err = NewAwsS3(bucket, basePath, region, "", accessKeyID, secretAccessKey)
	if err != nil {
		t.Log(err)
		return
	}
	fmt.Println(*as3)
}

func TestAwsS3_UploadFromFile(t *testing.T) {
	localFile := "README.md"
	utils.SafeRunWithTimeout(time.Second*2, func(cancel context.CancelFunc) {
		url, err := testS3.UploadFromFile(localFile)
		if err != nil {
			t.Log(err)
			return
		}
		fmt.Println("upload success, url =", url)
		cancel()
	})
}

func TestAwsS3_UploadFromReader(t *testing.T) {
	localFile := "./README.md"
	f, err := os.Open(localFile)
	if err != nil {
		t.Log(err)
		return
	}
	defer f.Close()

	utils.SafeRunWithTimeout(time.Second*2, func(cancel context.CancelFunc) {
		url, err := testS3.UploadFromReader(f, localFile)
		if err != nil {
			t.Log(err)
			return
		}
		fmt.Println("upload success, url =", url)
		cancel()
	})
}

func TestAwsS3_DownloadToFile(t *testing.T) {
	awsFile := "README.md"
	localFile := "./README.md"
	utils.SafeRunWithTimeout(time.Second*2, func(cancel context.CancelFunc) {
		n, err := testS3.DownloadToFile(awsFile, localFile)
		if err != nil {
			t.Log(err)
			return
		}
		fmt.Println("download file success, size =", n)
		cancel()
	})
}

func TestAwsS3_CheckFileIsExist(t *testing.T) {
	awsFile := "README.md"
	utils.SafeRunWithTimeout(time.Second*2, func(cancel context.CancelFunc) {
		err := testS3.CheckFileIsExist(awsFile)
		if err != nil {
			t.Log(err)
			return
		}
		fmt.Println(awsFile, "is exist")
		cancel()
	})
}

func TestAwsS3_DeleteFile(t *testing.T) {
	awsFile := "README.md"
	utils.SafeRunWithTimeout(time.Second*2, func(cancel context.CancelFunc) {
		err := testS3.DeleteFile(awsFile)
		if err != nil {
			t.Log(err)
			return
		}
		fmt.Printf("delete file %s success.\n", awsFile)
		cancel()
	})
}

func TestAwsS3_GetPreSignedURL(t *testing.T) {
	errMsg := ""
	awsFiles := []string{"README.md"}
	for _, awsFile := range awsFiles {
		utils.SafeRunWithTimeout(time.Second*2, func(cancel context.CancelFunc) {
			url, err := testS3.GetPreSignedURL(DownloadMethod, awsFile, 1000)
			if err != nil {
				errMsg += err.Error() + "\n"
				return
			}
			fmt.Println(url)
			cancel()
		})
	}

	if errMsg != "" {
		t.Log(errMsg)
	}
}

func TestAwsS3_WithBucketPath(t *testing.T) {
	currentS3 := testS3.WithBucketAndBasePath("/myPath/")
	fmt.Println(testS3.Bucket, testS3.BasePath)
	fmt.Println(currentS3.Bucket, currentS3.BasePath)
	fmt.Println(testS3.S3, testS3.Session)
	fmt.Println(currentS3.S3, currentS3.Session)

	utils.SafeRunWithTimeout(time.Second*2, func(cancel context.CancelFunc) {
		awsFile := "README.md"
		err := testS3.CheckFileIsExist(awsFile)
		if err != nil {
			t.Log(err)
		} else {
			fmt.Println(awsFile, "is exist")
		}
		cancel()
	})

	utils.SafeRunWithTimeout(time.Second*2, func(cancel context.CancelFunc) {
		awsFile := "README.md"
		err = currentS3.WithBucketAndBasePath("/2019/04/05/", "dev-k8s-log").CheckFileIsExist(awsFile)
		if err != nil {
			t.Log(err)
		} else {
			fmt.Println(awsFile, "is exist")
		}
		cancel()
	})
}
