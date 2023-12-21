package rules

import (
	"context"
)

func GetAll(ctx context.Context) ([]Rule, error) {

	/*gitHubClient, err := github.CreateGithubClient(ctx) //TODO
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred creating the GitHub client")
	}*/

	allRules := []Rule{
		newDuplicatedPackageRule(),
		//newValidPackageRule(gitHubClient), //TODO
		//newValidPackageIconRule(gitHubClient), //TODO
	}

	return allRules, nil
}
