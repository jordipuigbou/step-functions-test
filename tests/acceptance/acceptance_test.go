package main

import (
	"context"
	"os"
	"testing"

	"github.com/jordipuigbou/step-functions-test/tests/acceptance/steps/awss"
	"github.com/jordipuigbou/step-functions-test/tests/acceptance/steps/db"
	"github.com/jordipuigbou/step-functions-test/tests/acceptance/steps/iam"
	"github.com/jordipuigbou/step-functions-test/tests/acceptance/steps/lambda"
	"github.com/jordipuigbou/step-functions-test/tests/acceptance/steps/s3"
	"github.com/jordipuigbou/step-functions-test/tests/acceptance/steps/sfn"

	"github.com/TelefonicaTC2Tech/golium"
	"github.com/TelefonicaTC2Tech/golium/steps/common"
	"github.com/cucumber/godog"
)

func TestMain(m *testing.M) {
	launcher := golium.NewLauncher()
	launcher.Launch(InitializeTestSuite, InitializeScenario)
	exitVal := m.Run()
	os.Exit(exitVal)
}

func InitializeTestSuite(ctx context.Context, suiteCtx *godog.TestSuiteContext) {
}

func InitializeScenario(ctx context.Context, scenarioCtx *godog.ScenarioContext) {
	stepsInitializers := []golium.StepsInitializer{
		common.Steps{},
		sfn.Steps{},
		db.Steps{},
		lambda.Steps{},
		s3.Steps{},
		iam.Steps{},
		awss.Steps{},
	}
	for _, stepsInitializer := range stepsInitializers {
		ctx = stepsInitializer.InitializeSteps(ctx, scenarioCtx)
	}

	// Initialize dependencies
	awsFeature := &AwsFeature{}
	awsFeature.setSessions(ctx)
	awsFeature.propagateAwsSession(ctx)
	// Scenario Setup and Teardown
	awsFeature.beforeScenario()

	scenarioCtx.After(awsFeature.afterScenario)
}

type AwsFeature struct {
	SfnSession    *sfn.Session
	DBSession     *db.Session
	LambdaSession *lambda.Session
	S3Session     *s3.Session
	IAMSession    *iam.Session
	AWSSession    *awss.Session
}

func (a *AwsFeature) beforeScenario() {
	a.SfnSession.SetAwsSfnClient()
	a.IAMSession.SetIAMClient()
	a.DBSession.SetAwsDynamoClient()
	a.LambdaSession.SetLambdaClient()
	a.S3Session.SetS3Client()
	a.S3Session.SetS3UploaderManager()
}

func (a *AwsFeature) afterScenario(ctx context.Context, s *godog.Scenario, testErr error) (
	context.Context, error) {
	if err := a.IAMSession.DeleteIAMRole(ctx); err != nil {
		return ctx, err
	}
	if err := a.SfnSession.DeleteSfnStateMachine(); err != nil {
		return ctx, nil
	}
	if err := a.LambdaSession.DeleteLambdaFunction(); err != nil {
		return ctx, err
	}
	if err := a.DBSession.DeleteDynamoDBTable(a.DBSession.TestTableName); err != nil {
		return ctx, err
	}
	return ctx, nil
}

func (a *AwsFeature) setSessions(ctx context.Context) {
	a.AWSSession = awss.GetSession(ctx)
	a.SfnSession = sfn.GetSession(ctx)
	a.DBSession = db.GetSession(ctx)
	a.LambdaSession = lambda.GetSession(ctx)
	a.S3Session = s3.GetSession(ctx)
	a.LambdaSession.S3Session = a.S3Session
	a.IAMSession = iam.GetSession(ctx)
}
func (a *AwsFeature) propagateAwsSession(ctx context.Context) {
	a.AWSSession.SetAwsSession(ctx)
	a.SfnSession.AwsSession = a.AWSSession.AwsSession
	a.LambdaSession.AwsSession = a.AWSSession.AwsSession
	a.S3Session.AwsSession = a.AWSSession.AwsSession
	a.DBSession.AwsSession = a.AWSSession.AwsSession
	a.IAMSession.AwsSession = a.AWSSession.AwsSession
}
