package main

import (
	"context"
	"os"
	"testing"

	"github.com/jordipuigbou/step-functions-test/tests/acceptance/steps/db"
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
	}
	for _, stepsInitializer := range stepsInitializers {
		ctx = stepsInitializer.InitializeSteps(ctx, scenarioCtx)
	}

	// Initialize dependencies
	awsFeature := &AwsFeature{}
	awsFeature.SfnSession = sfn.GetSession(ctx)
	awsFeature.DBSession = db.GetSession(ctx)

	// Scenario Setaup and Teardown
	awsFeature.beforeScenario(ctx)

	scenarioCtx.After(awsFeature.afterScenario)
}

type AwsFeature struct {
	SfnSession *sfn.Session
	DBSession  *db.Session
}

func (a *AwsFeature) beforeScenario(ctx context.Context) (
	context.Context, error) {
	err := a.SfnSession.SetAwsSfnClient(ctx)
	if err != nil {
		return ctx, err
	}
	err = a.DBSession.SetAwsDynamoClient(ctx)
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}

func (a *AwsFeature) afterScenario(ctx context.Context, s *godog.Scenario, testErr error) (
	context.Context, error) {
	if err := a.SfnSession.DeleteIAMRole(ctx); err != nil {
		return ctx, err
	}
	if err := a.SfnSession.DeleteSfnStateMachine(); err != nil {
		return ctx, nil
	}

	if err := a.SfnSession.DeleteLambdaFunction(); err != nil {
		return ctx, err
	}
	if err := a.DBSession.DeleteDynamoDBTable(a.DBSession.TestTableName); err != nil {
		return ctx, err
	}
	return ctx, nil
}
