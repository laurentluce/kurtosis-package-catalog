package rules

import (
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/types"
	"github.com/kurtosis-tech/stacktrace"
)

type CheckResult struct {
	ruleName     RuleName
	wasValidated bool
	failures     map[types.PackageName][]string
}

func newCheckResult(ruleName RuleName, wasValidated bool, failures map[types.PackageName][]string) *CheckResult {
	return &CheckResult{ruleName: ruleName, wasValidated: wasValidated, failures: failures}
}

func (ruleReport *CheckResult) GetRuleName() RuleName {
	return ruleReport.ruleName
}

func (ruleReport *CheckResult) WasValidated() bool {
	return ruleReport.wasValidated
}

func (ruleReport *CheckResult) GetFailures() map[types.PackageName][]string {
	return ruleReport.failures
}

func (ruleReport *CheckResult) GetFailuresForPackage(packageName types.PackageName) ([]string, error) {
	failures, found := ruleReport.failures[packageName]
	if !found {
		return nil, stacktrace.NewError("Expected to find failures for package '%s' but nothing was found, this is a bug in the catalog", packageName)
	}
	return failures, nil
}
