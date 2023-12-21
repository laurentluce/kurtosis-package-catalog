package rules

import (
	"context"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/catalog"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/types"
	"github.com/kurtosis-tech/stacktrace"
)

const (
	duplicatedPackageRuleName = "Duplicated package"
)

// duplicatedPackageRule checks that there is not duplicated packages name in the catalog
type duplicatedPackageRule struct {
	name string
}

func newDuplicatedPackageRule() *duplicatedPackageRule {
	return &duplicatedPackageRule{name: duplicatedPackageRuleName}
}

func (duplicatedPackageRule *duplicatedPackageRule) GetName() string {
	return duplicatedPackageRule.name
}

func (duplicatedPackageRule *duplicatedPackageRule) Check(ctx context.Context, catalog catalog.PackageCatalog) error {

	packageNames := map[types.PackageName]bool{}

	for _, packageData := range catalog {
		packageName := packageData.GetPackageName()
		if _, found := packageNames[packageName]; found {
			return stacktrace.NewError("duplicated package name '%s' found in the Kurtosis package catalog", packageName)
		}
		packageNames[packageName] = true
	}
	return nil
}
