package validator

import (
	"context"
	"github.com/kurtosis-tech/kurtosis-package-catalog/catalog-validator/validation/rules"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/catalog"
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

func (validator *Validator) Validate(ctx context.Context) error {
	for _, rule := range validator.rules {
		logrus.Debugf("Checking rule '%s'", rule.GetName())
		if err := rule.Check(ctx, validator.catalog); err != nil {
			return stacktrace.Propagate(err, "invalid Kurtosis package catalog, the current version do not pass the '%s' validation rule", rule.GetName())
		}
		logrus.Debugf("'%s' rule passed", rule.GetName())
	}
	return nil
}
