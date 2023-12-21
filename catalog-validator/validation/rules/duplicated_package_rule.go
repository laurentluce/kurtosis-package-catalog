package rules

import (
	"context"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/catalog"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/types"
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

func (duplicatedPackageRule *duplicatedPackageRule) GetName() RuleName {
	return RuleName(duplicatedPackageRule.name)
}

func (duplicatedPackageRule *duplicatedPackageRule) Check(_ context.Context, catalog catalog.PackageCatalog) *CheckResult {

	wasValidated := true
	failures := map[types.PackageName][]string{}

	packageNames := map[types.PackageName]bool{}

	for _, packageData := range catalog {
		packageName := packageData.GetPackageName()
		if _, found := packageNames[packageName]; found {
			failures[packageName] = []string{"duplicated name"}
			wasValidated = false
			continue
		}
		packageNames[packageName] = true
	}

	checkResult := newCheckResult(duplicatedPackageRule.GetName(), wasValidated, failures)

	return checkResult
}
