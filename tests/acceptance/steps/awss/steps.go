package awss

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
	scenCtx.Step(`^I set AWS Session$`,
		func() error {
			return session.SetAwsSession(ctx)
		})
	return ctx
}
