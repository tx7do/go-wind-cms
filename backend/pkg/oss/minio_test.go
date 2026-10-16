package oss

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
	"github.com/tx7do/go-utils/trans"

	storageV1 "go-wind-cms/api/gen/go/storage/service/v1"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func createTestClient() *MinIOClient {
	return NewMinIoClient(&conf.Bootstrap{
		Oss: &conf.OSS{
			Minio: &conf.OSS_MinIO{
				Endpoint:     "127.0.0.1:9000",
				UploadHost:   "127.0.0.1:9000",
				DownloadHost: "127.0.0.1:9000",
				AccessKey:    "root",
				SecretKey:    "*Abcd123456",
			},
		},
	}, log.DefaultLogger)
}

func TestMinIoClient(t *testing.T) {
	cli := createTestClient()
	assert.NotNil(t, cli)

	resp, err := cli.GetUploadPresignedUrl(t.Context(), &storageV1.GetUploadPresignedUrlRequest{
		Method:        storageV1.GetUploadPresignedUrlRequest_Put,
		ContentType:   trans.String("image/jpeg"),
		BucketName:    trans.String("images"),
		FileDirectory: trans.String("20221010"),
	})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestListFile(t *testing.T) {
	cli := createTestClient()
	assert.NotNil(t, cli)

	req := &storageV1.ListOssFileRequest{
		BucketName: trans.Ptr("users"),
		Folder:     trans.Ptr("1"),
		Recursive:  trans.Ptr(true),
	}
	files, err := cli.ListFile(t.Context(), req)
	assert.Nil(t, err)
	fmt.Println(files)
}

func TestDownloadFile(t *testing.T) {
	cli := createTestClient()
	assert.NotNil(t, cli)

	// 自播种：下载测试不再依赖外部 MinIO 里恰好存在的对象（fixture 早已漂移），
	// 先上传唯一对象再下载，测试自清理。
	bucket := "images"
	object := fmt.Sprintf("test/download-seed-%d.png", time.Now().UnixNano())
	payload := []byte("download-seed-payload")
	_, _, _, err := cli.UploadFile(t.Context(), bucket, object, "image/png",
		bytes.NewReader(payload), int64(len(payload)))
	assert.Nil(t, err)
	defer func() {
		_ = cli.DeleteFile(t.Context(), bucket, object)
	}()

	resp, derr := cli.DownloadFile(t.Context(), &storageV1.DownloadFileRequest{
		Selector: &storageV1.DownloadFileRequest_StorageObject{
			StorageObject: &storageV1.StorageObject{
				BucketName: trans.Ptr(bucket),
				ObjectName: trans.Ptr(object),
			},
		},
		PreferPresignedUrl: trans.Ptr(false),
	})
	assert.Nil(t, derr)
	assert.NotNil(t, resp)
}
