package iam

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	aws_s "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

const (
	role       = "admin"
	rolePolicy = `{"Version":"2012-10-17","Statement":[{"Effect":"Allow"` +
		`,"Action":"*","Resource":"*"}]}`
)

type Session struct {
	AwsSession *aws_s.Session
	iamClient  *iam.IAM
}

func (s *Session) SetIAMClient() {
	s.iamClient = iam.New(s.AwsSession)
}
func (s *Session) CreateIAMRole(ctx context.Context) error {
	logger := GetLogger()
	logger.LogMessage("Creating " + role + " role...")
	iamRoleOut, err := s.iamClient.CreateRole(&iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(rolePolicy),
		RoleName:                 aws.String(role),
		Description:              aws.String("Testing step function and lambda role"),
	})
	if err != nil {
		return fmt.Errorf("error creating iam role: %v", err)
	}
	logger.LogMessage("new " + role + " iam role created with...")
	logger.LogMessage(*iamRoleOut.Role.RoleName)
	logger.LogMessage(*iamRoleOut.Role.Description)
	logger.LogMessage(*iamRoleOut.Role.RoleId)
	logger.LogMessage(*iamRoleOut.Role.Arn)

	return nil
}

func (s *Session) DeleteIAMRole(ctx context.Context) error {
	logger := GetLogger()
	logger.LogMessage("Deleting IAM " + role + " role...")
	_, err := s.iamClient.DeleteRole(&iam.DeleteRoleInput{
		RoleName: aws.String(role),
	})

	if err != nil {
		return fmt.Errorf("error deleting iam role: %v", err)
	}
	logger.LogMessage("IAM " + role + " deleted")
	return nil
}
