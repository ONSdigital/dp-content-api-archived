package s3

import (
	"io/ioutil"
	"log"

	"github.com/minio/minio-go"
	"github.com/minio/minio-go/pkg/credentials"
)

type MinioClient struct {
	client *minio.Client
	bucket string
}

func (s3 *MinioClient) GetBucket() string {
	return s3.bucket
}

func CreateClient(awsZone, bucket, endpoint string, iam, useSSL bool) S3Client {
	var cred *credentials.Credentials
	if iam {
		cred = credentials.NewIAM("")
	} else {
		cred = credentials.NewFileMinioClient("", "s3")
	}
	minioClient, err := minio.NewWithCredentials(endpoint, cred, useSSL, awsZone)
	if err != nil {
		log.Fatalln(err)
	}
	return &MinioClient{minioClient, bucket}
}

func (s3 *MinioClient) GetReader(location string) (*minio.Object, error) {
	object, err := s3.client.GetObject(s3.bucket, location)
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (s3 *MinioClient) GetObject(location string) ([]byte, error) {
	object, err := s3.GetReader(location)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(object)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s3 *MinioClient) GetFileSize(location string) (int64, error) {
	object, err := s3.client.GetObject(s3.bucket, location)
	if err != nil {
		return 0, err
	}
	stat, err := object.Stat()
	if err != nil {
		return 0, err
	}
	return stat.Size, nil
}
