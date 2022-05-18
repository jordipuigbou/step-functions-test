package awss

import (
	"context"
	"fmt"

	"github.com/TelefonicaTC2Tech/golium"
	"github.com/aws/aws-sdk-go/aws"
	aws_s "github.com/aws/aws-sdk-go/aws/session"
)

const (
	role       = "admin"
	rolePolicy = `{"Version":"2012-10-17","Statement":[{"Effect":"Allow"` +
		`,"Action":"*","Resource":"*"}]}`
)

type Session struct {
	AwsSession *aws_s.Session
}

func (s *Session) SetAwsSession(ctx context.Context) error {
	awsConfig := &aws.Config{
		Endpoint:                      aws.String(golium.Value(ctx, "[CONF:awsEndpoint]").(string)),
		DisableSSL:                    aws.Bool(true),
		Region:                        aws.String(golium.Value(ctx, "[CONF:awsRegion]").(string)),
		DisableEndpointHostPrefix:     aws.Bool(true),
		S3ForcePathStyle:              aws.Bool(true),
		CredentialsChainVerboseErrors: aws.Bool(true),
	}
	var err error
	if s.AwsSession, err = aws_s.NewSession(awsConfig); err != nil {
		return fmt.Errorf("error creating aws session. %v", err)
	}
	return nil
}
