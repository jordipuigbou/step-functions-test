package iam

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
	scenCtx.Step(`^I create iam testing role$`,
		func() error {
			return session.CreateIAMRole(ctx)
		})
	return ctx
}
