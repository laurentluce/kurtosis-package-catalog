package rules

import (
	"context"
	"fmt"
	"github.com/google/go-github/v54/github"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/catalog"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/consts"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/types"
	"github.com/kurtosis-tech/stacktrace"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"net/http"
	"path"
)

const (
	validPackageRuleName = "Valid package"
)

type KurtosisYaml struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// validPackageRule checks if the package is valid by checking if:
// 1- the package repository exist
// 2- if the package repository contains the kurtosis.yml file
// 3- if the name inside the kurtosis.yml file is the same in the package catalog
type validPackageRule struct {
	name         string
	gitHubClient *github.Client
}

func newValidPackageRule(gitHubClient *github.Client) *validPackageRule {
	return &validPackageRule{name: validPackageRuleName, gitHubClient: gitHubClient}
}

func (validPackageRule *validPackageRule) GetName() RuleName {
	return RuleName(validPackageRule.name)
}

func (validPackageRule *validPackageRule) Check(ctx context.Context, catalog catalog.PackageCatalog) *CheckResult {

	wasValidated := true
	failures := map[types.PackageName][]string{}

	for _, packageData := range catalog {
		packageName := packageData.GetPackageName()
		logrus.Debugf("Checking if package '%s' is valid...", packageName)
		repositoryOwner := packageData.GetRepositoryOwner()
		repositoryName := packageData.GetRepositoryName()
		repositoryPackageRootPath := packageData.GetRepositoryPackageRootPath()
		packageFailures := []string{}
		packageNameFromKurtosisYamlFile, err := validPackageRule.getPackageNameFromKurtosisYmlFile(ctx, packageName, repositoryOwner, repositoryName, repositoryPackageRootPath)
		if err != nil {
			errorFailure := fmt.Sprintf("the package does not exist or does not contains the '%s' file", consts.DefaultKurtosisYamlFilename)
			packageFailures = append(packageFailures, errorFailure)
		} else {
			if packageName != packageNameFromKurtosisYamlFile {
				invalidPackageNameMsg := fmt.Sprintf("package name '%s' in the catalog does not match with the name '%s' found in the package repository", packageName, packageNameFromKurtosisYamlFile)
				packageFailures = append(packageFailures, invalidPackageNameMsg)
			}
		}

		if len(packageFailures) > 0 {
			failures[packageName] = packageFailures
			wasValidated = false
			continue
		}
		logrus.Debugf("...package '%s' successfully validated.", packageName)
	}

	checkResult := newCheckResult(validPackageRule.GetName(), wasValidated, failures)

	return checkResult
}

func (validPackageRule *validPackageRule) getPackageNameFromKurtosisYmlFile(ctx context.Context, packageName types.PackageName, repositoryOwner string, repositoryName string, repositoryPackageRootPath string) (types.PackageName, error) {
	kurtosisYamlFilepath := path.Join(repositoryPackageRootPath, consts.DefaultKurtosisYamlFilename)
	repoGetContentOpts := &github.RepositoryContentGetOptions{
		Ref: "",
	}

	// get contents of kurtosis yaml file from GitHub
	kurtosisYamlFileContentResult, _, resp, err := validPackageRule.gitHubClient.Repositories.GetContents(ctx, repositoryOwner, repositoryName, kurtosisYamlFilepath, repoGetContentOpts)
	if err != nil && resp != nil && resp.StatusCode == http.StatusNotFound {
		return "", stacktrace.NewError("No '%s' file for package '%s'", kurtosisYamlFilepath, packageName)
	} else if err != nil {
		return "", stacktrace.Propagate(err, "An error occurred reading content of Kurtosis Package '%s' - file '%s'", packageName, kurtosisYamlFilepath)
	}

	kurtosisYaml, err := parseKurtosisYaml(kurtosisYamlFileContentResult)
	if err != nil {
		return "", stacktrace.Propagate(err, "an error occurred parsing the Kurtosis YAML file for '%s'", packageName)
	}

	packageNameFromYamlFile := types.PackageName(kurtosisYaml.Name)

	return packageNameFromYamlFile, nil
}

func parseKurtosisYaml(kurtosisYamlContent *github.RepositoryContent) (*KurtosisYaml, error) {
	rawFileContent, err := kurtosisYamlContent.GetContent()
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred getting the content of the '%s' file", consts.DefaultKurtosisYamlFilename)
	}

	kurtosisYaml := new(KurtosisYaml)
	if err = yaml.Unmarshal([]byte(rawFileContent), kurtosisYaml); err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred parsing YAML for '%s'", consts.DefaultKurtosisYamlFilename)
	}

	if kurtosisYaml.Name == "" {
		return nil, stacktrace.NewError("Kurtosis YAML file had an empty name. This is invalid.")
	}
	return kurtosisYaml, nil
}
