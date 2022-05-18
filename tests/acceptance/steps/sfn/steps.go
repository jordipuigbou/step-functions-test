package sfn

import (
	"context"

	"github.com/TelefonicaTC2Tech/golium"
	"github.com/cucumber/godog"
)

type Steps struct {
}

func (s Steps) InitializeSteps(
	ctx context.Context,
	scenCtx *godog.ScenarioContext,
) context.Context {
	ctx = InitializeContext(ctx)
	session := GetSession(ctx)
	// Initialize the steps
	scenCtx.Step(`^I set AWS SFN client$`,
		func() {
			session.SetAwsSfnClient()
		})
	scenCtx.Step(`^I create step function state machine from "([^"]*)" JSON file with "([^"]*)" name$`,
		func(stepFunction, stateMachineName string) error {
			return session.CreateStepFunctionMachine(ctx, stateMachineName, stepFunction)
		})
	scenCtx.Step(`^I start step function execution with "([^"]*)" name$`,
		func(executionName string) error {
			return session.StartSfnExecution(executionName)
		})
	scenCtx.Step(`^I get step function execution describe and store response in "([^"]*)" key$`,
		func(key string) error {
			return session.DescribeSfnStateExecution(ctx, key)
		})
	scenCtx.Step(`^I validate step function output stored in "([^"]*)" key is equal to$`,
		func(key string, message *godog.DocString) error {
			return session.ValidateSfnOutput(ctx, key, golium.ValueAsString(ctx, message.Content))
		})

	scenCtx.Step(`^I delete step function state machine$`,
		func() error {
			return session.DeleteSfnStateMachine()
		})

	return ctx
}
