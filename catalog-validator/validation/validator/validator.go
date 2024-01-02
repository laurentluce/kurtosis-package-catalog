package validator

import (
	"context"
	"github.com/kurtosis-tech/kurtosis-package-catalog/catalog-validator/validation/rules"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/catalog"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/types"
	"github.com/kurtosis-tech/stacktrace"
	"github.com/sirupsen/logrus"
)

type Validator struct {
	catalog catalog.PackageCatalog
	rules   []rules.Rule
}

func NewValidator(catalog catalog.PackageCatalog, rules []rules.Rule) *Validator {
	return &Validator{catalog: catalog, rules: rules}
}

func (validator *Validator) Validate(ctx context.Context) (*result, error) {

	isValidCatalog := true
	rulesResult := map[rules.RuleName]map[types.PackageName][]string{}

	for _, rule := range validator.rules {
		logrus.Debugf("Checking rule '%s'", rule.GetName())
		if checkResult := rule.Check(ctx, validator.catalog); !checkResult.WasValidated() {
			isValidCatalog = false
			failures := checkResult.GetFailures()
			var packageName types.PackageName
			for packageNameInFailures := range failures {
				packageName = packageNameInFailures
				break
			}
			ruleName := checkResult.GetRuleName()
			failuresByPackageForRule, found := rulesResult[ruleName]
			if found {
				packageFailures, err := checkResult.GetFailuresForPackage(packageName)
				if err != nil {
					return nil, stacktrace.Propagate(err, "an error occurred getting failures for package '%s'", packageName)
				}
				failuresByPackageForRule[packageName] = packageFailures
			} else {
				rulesResult[ruleName] = checkResult.GetFailures()
			}

			logrus.Debugf("the current catalog version does not pass rule '%s'", rule.GetName())
			continue
		}
		logrus.Debugf("'%s' rule passed", rule.GetName())
	}

	resultObj := newResult(isValidCatalog, rulesResult)

	return resultObj, nil
}
