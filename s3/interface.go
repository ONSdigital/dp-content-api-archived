package s3

type S3Client interface {
	GetObject(string) ([]byte, error)
	GetBucket() string
	GetFileSize(string) (int64, error)
}
