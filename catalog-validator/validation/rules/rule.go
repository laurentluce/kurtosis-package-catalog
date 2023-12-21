package rules

import (
	"context"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/catalog"
)

type RuleName string

type Rule interface {
	GetName() RuleName
	Check(ctx context.Context, catalog catalog.PackageCatalog) *CheckResult
}
