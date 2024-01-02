package validator

import (
	"github.com/kurtosis-tech/kurtosis-package-catalog/catalog-validator/validation/rules"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/types"
)

type result struct {
	isValidCatalog bool
	rulesResult    map[rules.RuleName]map[types.PackageName][]string
}

func newResult(isValidCatalog bool, rulesResult map[rules.RuleName]map[types.PackageName][]string) *result {
	return &result{isValidCatalog: isValidCatalog, rulesResult: rulesResult}
}

func (result *result) IsValidCatalog() bool {
	return result.isValidCatalog
}

func (result *result) GetRulesResult() map[rules.RuleName]map[types.PackageName][]string {
	return result.rulesResult
}
