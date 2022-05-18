package lambda

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	aws_s "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/jordipuigbou/step-functions-test/tests/acceptance/steps/s3"
)

const (
	projectPath           = "../../"
	lambdaFunctionName    = "HelloWorldFunction"
	lambdaFunctionPath    = projectPath + "/.aws-sam/build"
	lambdaFunctionZipFile = lambdaFunctionName + ".zip"
	bucket                = "bucketname"
)

type Session struct {
	AwsSession   *aws_s.Session
	S3Session    *s3.Session
	lambdaClient *lambda.Lambda
}

func (s *Session) SetLambdaClient() {
	s.lambdaClient = lambda.New(s.AwsSession)
}

func (s *Session) CreateLambdaFunction(ctx context.Context) error {
	logger := GetLogger()
	logger.LogMessage("Creating lambda function...")

	if err := s.S3Session.CreateS3Bucket(bucket); err != nil {
		return fmt.Errorf("error creating s3 bucket: %v", err)
	}

	if err := s.S3Session.UploadFileToS3Bucket(
		lambdaFunctionPath, lambdaFunctionZipFile, bucket,
	); err != nil {
		return fmt.Errorf("error uploading zip file to s3 bucket: %v", err)
	}

	lambdaOut, err := s.lambdaClient.CreateFunction(&lambda.CreateFunctionInput{
		FunctionName: aws.String(lambdaFunctionName),
		Runtime:      aws.String("go1.x"),
		Role:         aws.String("arn:aws:iam::000000000000:role/admin"),
		Handler:      aws.String("hello_world"),
		Code: &lambda.FunctionCode{
			S3Bucket: aws.String(bucket),
			S3Key:    aws.String(lambdaFunctionZipFile),
		},
	})
	logger.LogMessage("Create lambda function finished")
	if err != nil {
		return fmt.Errorf("error creating lambda function: %v", err)
	}
	logger.LogMessage("Created lambda function from S3 bucket...")
	logger.LogMessage(*lambdaOut.FunctionName)
	logger.LogMessage(*lambdaOut.FunctionArn)
	logger.LogMessage(*lambdaOut.Handler)
	logger.LogMessage(*lambdaOut.Role)
	return nil
}

func (s *Session) DeleteLambdaFunction() error {
	logger := GetLogger()
	logger.LogMessage("Deleting lambda function")
	_, err := s.lambdaClient.DeleteFunction(&lambda.DeleteFunctionInput{
		FunctionName: aws.String(lambdaFunctionName),
	})
	if err != nil {
		return fmt.Errorf("error deleting lambda function: %v", err)
	}
	logger.LogMessage("Deleted lambda function")
	return nil
}
