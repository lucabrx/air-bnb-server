package aws

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type AWS struct {
	AccessKey string
	SecretKey string
}

func NewAws(accessKey, secretKey string) AWS {
	return AWS{
		SecretKey: secretKey,
		AccessKey: accessKey,
	}
}

func (a *AWS) Client() (*minio.Client, error) {
	s3Client, err := minio.New("s3.eu-central-1.amazonaws.com", &minio.Options{
		Creds:  credentials.NewStaticV4(a.AccessKey, a.SecretKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, err
	}

	return s3Client, nil
}
