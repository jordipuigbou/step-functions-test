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
	"github.com/aws/aws-sdk-go/service/sfn"
)

const (
	projectPath       = "../../"
	stateMachinesPath = "statemachines"
)

type Session struct {
	AwsSession    *aws_s.Session
	sfnClient     *sfn.SFN
	stateMachine  *sfn.CreateStateMachineOutput
	lastExecution *sfn.StartExecutionOutput
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

func (s *Session) SetAwsSfnClient() {
	s.sfnClient = sfn.New(s.AwsSession)
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
