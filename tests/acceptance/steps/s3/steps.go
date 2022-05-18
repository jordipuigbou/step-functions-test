package s3

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
	_ = GetSession(ctx)
	// Initialize the steps

	return ctx
}
