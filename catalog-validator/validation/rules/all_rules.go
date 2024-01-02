package rules

import (
	"context"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/github"
	"github.com/kurtosis-tech/stacktrace"
)

func GetAll(ctx context.Context) ([]Rule, error) {

	gitHubClient, err := github.CreateGithubClient(ctx)
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred creating the GitHub client")
	}

	allRules := []Rule{
		newDuplicatedPackageRule(),
		newValidPackageRule(gitHubClient),
		newValidPackageIconRule(gitHubClient),
	}

	return allRules, nil
}
