package db

import (
	"context"

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
	scenCtx.Step(`^I set AWS DynamoDB client$`,
		func() {
			session.SetAwsDynamoClient()
		})
	scenCtx.Step(`^I create a DynamoDB table with "([^"]*)" name and "([^"]*)" string index$`,
		func(tableName, indexName string) error {
			return session.CreateDynamoDBTable(tableName, indexName)
		})

	scenCtx.Step(`^I delete a DynamoDB table with "([^"]*)" name$`,
		func(tableName string) error {
			return session.DeleteDynamoDBTable(tableName)
		})

	return ctx
}
