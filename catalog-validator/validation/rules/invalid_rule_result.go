package rules

import "github.com/kurtosis-tech/kurtosis-package-indexer/server/types"

type CheckResult struct {
	ruleName     RuleName
	wasValidated bool //TODO look for a better name
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

func (ruleReport *CheckResult) GetFailuresForPackage(packageName types.PackageName) []string {
	return ruleReport.failures[packageName] // TODO check if found
}
