package rules

import (
	"context"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/catalog"
)

type Rule interface {
	GetName() string
	Check(ctx context.Context, catalog catalog.PackageCatalog) error
}
