package s3

import (
	"fmt"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	aws_s "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	bucket = "bucketname"
)

type Session struct {
	AwsSession        *aws_s.Session
	s3Client          *s3.S3
	s3UploaderManager *s3manager.Uploader
}

func (s *Session) SetS3Client() {
	s.s3Client = s3.New(s.AwsSession)
}

func (s *Session) SetS3UploaderManager() {
	s.s3UploaderManager = s3manager.NewUploader(s.AwsSession)
}
func (s *Session) CreateS3Bucket(bucketName string) error {

	if _, err := s.s3Client.CreateBucket(&s3.CreateBucketInput{
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String("us-east-1"),
		},
		Bucket: aws.String(bucketName),
	}); err != nil {
		return fmt.Errorf("error creating a new bucket: %s, err: %v", bucket, err)
	}
	return nil
}

func (s *Session) UploadFileToS3Bucket(filePath, fileName, bucketName string) error {
	file, err := os.Open(path.Join(filePath, fileName))
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", fileName, err)
	}
	defer file.Close()

	_, err = s.s3UploaderManager.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("error uploading zip file to S3: %v", err)
	}
	return nil
}
