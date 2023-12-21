package validator

import (
	"github.com/kurtosis-tech/kurtosis-package-catalog/catalog-validator/validation/rules"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/types"
)

type result struct {
	isValidCatalog bool
	rulesResult    map[types.PackageName]map[rules.RuleName][]string
}

func newResult(isValidCatalog bool, rulesResult map[types.PackageName]map[rules.RuleName][]string) *result {
	return &result{isValidCatalog: isValidCatalog, rulesResult: rulesResult}
}

func (result *result) IsValidCatalog() bool {
	return result.isValidCatalog
}

func (result *result) GetRulesResult() map[types.PackageName]map[rules.RuleName][]string {
	return result.rulesResult
}
