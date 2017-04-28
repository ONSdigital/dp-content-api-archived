package s3

import (
	"io/ioutil"
	"log"

	"github.com/minio/minio-go"
)

type MinioClient struct {
	client *minio.Client
	bucket string
}

func (s3 *MinioClient) GetBucket() string {
	return s3.bucket
}

func CreateClient(bucket, endpoint, accessKeyID, secretAccessKey string, useSSL bool) S3Client {
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
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
	data, _ := ioutil.ReadAll(object)
	if err != nil {
		return nil, err
	}
	return data, nil
}
