package uploader

import (
	"bytes"
	"sync"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/sirupsen/logrus"

	"bbs-go/common/urls"
	"bbs-go/config"
)

// 阿里云oss
type aliyunOssUploader struct {
	once   sync.Once
	bucket *oss.Bucket
}

func (aliyun *aliyunOssUploader) PutImage(data []byte) (string, error) {
	key := generateImageKey(data)
	return aliyun.PutObject(key, data)
}

func (aliyun *aliyunOssUploader) PutObject(key string, data []byte) (string, error) {
	bucket := aliyun.getBucket()
	if err := bucket.PutObject(key, bytes.NewReader(data)); err != nil {
		return "", err
	}
	c := config.Instance.Uploader.ObjectStorage
	return urls.UrlJoin(c.Host, key), nil
}

func (aliyun *aliyunOssUploader) CopyImage(originUrl string) (string, error) {
	data, err := download(originUrl)
	if err != nil {
		return "", err
	}
	return aliyun.PutImage(data)
}

func (aliyun *aliyunOssUploader) getBucket() *oss.Bucket {
	aliyun.once.Do(func() {
		c := config.Instance.Uploader.ObjectStorage
		if client, err := oss.New(c.Endpoint, c.AccessId, c.AccessSecret); err != nil {
			logrus.Error(err)
		} else if aliyun.bucket, err = client.Bucket(c.Bucket); err != nil {
			logrus.Error(err)
		}
	})
	return aliyun.bucket
}
