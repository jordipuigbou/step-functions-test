package sfn

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/TelefonicaTC2Tech/golium"
	"github.com/aws/aws-sdk-go/aws"
	aws_s "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sfn"
)

const (
	projectPath           = "../../"
	stateMachinesPath     = "statemachines"
	lambdaFunctionName    = "HelloWorldFunction"
	lambdaFunctionPath    = projectPath + "/.aws-sam/build"
	lambdaFunctionZipFile = lambdaFunctionName + ".zip"
	bucket                = "bucketname"
	role                  = "admin"
	rolePolicy            = `{"Version":"2012-10-17","Statement":[{"Effect":"Allow"` +
		`,"Action":"*","Resource":"*"}]}`
)

type Session struct {
	awsSession    *aws_s.Session
	sfnClient     *sfn.SFN
	stateMachine  *sfn.CreateStateMachineOutput
	lastExecution *sfn.StartExecutionOutput
	lambdaClient  *lambda.Lambda
}

type ConsumedCapacity struct {
	TableName     string  `json:"TableName"`
	CapacityUnits float64 `json:"CapacityUnits"`
}
type passResult struct {
	ConsumedCapacity ConsumedCapacity `json:"ConsumedCapacity"`
}

type executionOut struct {
	PassResult passResult `json:"passResult"`
	Message    string     `json:"message"`
}

func (s *Session) SetAwsSfnClient(ctx context.Context) error {
	awsConfig := &aws.Config{
		Endpoint:                  aws.String(golium.Value(ctx, "[CONF:sfnEndpoint]").(string)),
		DisableSSL:                aws.Bool(true),
		Region:                    aws.String(golium.Value(ctx, "[CONF:awsRegion]").(string)),
		DisableEndpointHostPrefix: aws.Bool(true),
		S3ForcePathStyle:          aws.Bool(true),
	}
	var err error
	if s.awsSession, err = aws_s.NewSession(awsConfig); err != nil {
		return fmt.Errorf("error creating aws session. %v", err)
	}
	s.sfnClient = sfn.New(s.awsSession)
	return nil
}

func (s *Session) CreateStepFunctionMachine(
	ctx context.Context, sfnStateMachineName, sfnDefinition string,
) error {
	logger := GetLogger()

	stepFunctionWithPath := path.Join(projectPath, stateMachinesPath, sfnDefinition)
	stepFunctionWithPath += ".json"

	cmd := exec.Command("ls", projectPath)
	stdout, _ := cmd.Output()

	// Print the output
	logger.LogMessage(string(stdout))

	definition, err := os.ReadFile(stepFunctionWithPath)
	if err != nil {
		return fmt.Errorf("error opening step function definition file: %v", err)
	}

	logger.LogMessage(string(definition))
	if s.sfnClient == nil {
		return fmt.Errorf("sfn client is not set")
	}
	var createdStateMachine *sfn.CreateStateMachineOutput
	start := time.Now()
	for {
		createdStateMachine, err = s.sfnClient.CreateStateMachine(&sfn.CreateStateMachineInput{
			Definition: aws.String(string(definition)),
			Name:       aws.String(sfnStateMachineName),
			RoleArn:    aws.String("arn:aws:iam::000000000000:role/admin"),
		})
		if err == nil || time.Since(start) > 30*time.Second {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return fmt.Errorf("error creating state machine: %v", err)
	}
	s.stateMachine = createdStateMachine
	return nil
}

func (s *Session) StartSfnExecution(executionName string) error {
	logger := GetLogger()
	if s.sfnClient == nil {
		return fmt.Errorf("sfn client is not set")
	}
	if s.stateMachine == nil {
		return fmt.Errorf("sfn state machine is not set")
	}
	var err error
	s.lastExecution, err = s.sfnClient.StartExecution(&sfn.StartExecutionInput{
		Name:            aws.String(executionName),
		StateMachineArn: s.stateMachine.StateMachineArn,
	})

	if err != nil {
		return fmt.Errorf("error starting execution: %v", err)
	}
	logger.LogMessage(*s.lastExecution.ExecutionArn)
	return nil
}

func (s *Session) DescribeSfnStateExecution(ctx context.Context, key string) error {
	logger := GetLogger()
	if s.sfnClient == nil {
		return fmt.Errorf("sfn client is not set")
	}
	if s.stateMachine == nil {
		return fmt.Errorf("sfn state machine is not set")
	}
	var executionOutput *sfn.DescribeExecutionOutput
	for {
		executionOutput, _ = checkExecution(s)
		logger.LogMessage(*executionOutput.Status)
		if *executionOutput.Status == "SUCCEEDED" {
			break
		}
		if *executionOutput.Status == "FAILED" {
			return fmt.Errorf("sfn task execution failed")
		}
		time.Sleep(1 * time.Second)
	}

	if executionOutput != nil {
		logger.LogMessage(*executionOutput.Output)
		logger.LogMessage(*executionOutput.Status)
	}

	output := executionOut{}
	if err := json.Unmarshal([]byte(*executionOutput.Output), &output); err != nil {
		return fmt.Errorf("unmarshalling error: %v", err)
	}
	logger.LogMessage(output.Message)
	logger.LogMessage(output.PassResult.ConsumedCapacity.TableName)
	golium.GetContext(ctx).Put(key, output.Message)
	return nil
}

func checkExecution(s *Session) (*sfn.DescribeExecutionOutput, error) {
	return s.sfnClient.DescribeExecution(&sfn.DescribeExecutionInput{
		ExecutionArn: s.lastExecution.ExecutionArn,
	})
}

func (s *Session) ValidateSfnOutput(ctx context.Context, key, expectedOutput string) error {
	output := golium.GetContext(ctx).Get(key)
	if output != expectedOutput {
		return fmt.Errorf("sfn output is not the expected: %v!=%v", output, expectedOutput)
	}
	return nil
}

func (s *Session) DeleteSfnStateMachine() error {
	logger := GetLogger()
	if s.sfnClient == nil {
		return fmt.Errorf("sfn client is not set")
	}
	if s.stateMachine == nil {
		return fmt.Errorf("sfn state machine is not set")
	}
	logger.LogMessage("Deleting sfn state machine...")
	_, err := s.sfnClient.DeleteStateMachine(&sfn.DeleteStateMachineInput{
		StateMachineArn: aws.String(*s.stateMachine.StateMachineArn),
	})
	if err != nil {
		return fmt.Errorf("error deleting state machine: %v", err)
	}
	logger.LogMessage("sfn state machine deleted")
	return nil
}

func (s *Session) CreateIAMRole(ctx context.Context) error {
	logger := GetLogger()
	iamClient := iam.New(s.awsSession)
	iamRoleOut, err := iamClient.CreateRole(&iam.CreateRoleInput{
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
	iamClient := iam.New(s.awsSession)
	_, err := iamClient.DeleteRole(&iam.DeleteRoleInput{
		RoleName: aws.String(role),
	})

	if err != nil {
		return fmt.Errorf("error deleting iam role: %v", err)
	}

	return nil
}

func (s *Session) CreateLambdaFunction(ctx context.Context) error {
	logger := GetLogger()
	logger.LogMessage("Creating lambda function...")
	s.lambdaClient = lambda.New(s.awsSession)

	if err := s.CreateS3Bucket(bucket); err != nil {
		return fmt.Errorf("error creating s3 bucket: %v", err)
	}

	if err := s.UploadFileToS3Bucket(lambdaFunctionPath, lambdaFunctionZipFile, bucket); err != nil {
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

func (s *Session) CreateS3Bucket(bucketName string) error {
	s3Client := s3.New(s.awsSession)
	if _, err := s3Client.CreateBucket(&s3.CreateBucketInput{
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

	uploader := s3manager.NewUploader(s.awsSession)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("error uploading zip file to S3: %v", err)
	}
	return nil
}

func (s *Session) DeleteLambdaFunction() error {
	logger := GetLogger()
	s.lambdaClient = lambda.New(s.awsSession)
	if s.lambdaClient == nil {
		return fmt.Errorf("lambda client is nil")
	}
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
