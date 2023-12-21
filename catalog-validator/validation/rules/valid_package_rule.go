package rules

import (
	"context"
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

func (packageExistRule *validPackageRule) GetName() string {
	return packageExistRule.name
}

func (packageExistRule *validPackageRule) Check(ctx context.Context, catalog catalog.PackageCatalog) error {

	for _, packageData := range catalog {
		packageName := packageData.GetPackageName()
		logrus.Debugf("Checking if package '%s' is valid...", packageName)
		repositoryOwner := packageData.GetRepositoryOwner()
		repositoryName := packageData.GetRepositoryName()
		repositoryPackageRootPath := packageData.GetRepositoryPackageRootPath()
		packageNameFromKurtosisYamlFile, err := packageExistRule.getPackageNameFromKurtosisYmlFile(ctx, packageName, repositoryOwner, repositoryName, repositoryPackageRootPath)
		if err != nil {
			return stacktrace.Propagate(err, "an error occurred getting the Kurtosis package name from the Kurtosis YAML file for package '%s'", packageName)
		}
		if packageName != packageNameFromKurtosisYamlFile {
			return stacktrace.NewError("there is an inconsistency between the Kurtosis package name in the catalog '%s' with the name '%s' found in the '%s' file ", packageName, packageNameFromKurtosisYamlFile, consts.DefaultKurtosisYamlFilename)
		}
		logrus.Debugf("...package '%s' successfully validated.", packageName)
	}

	return nil
}

func (packageExistRule *validPackageRule) getPackageNameFromKurtosisYmlFile(ctx context.Context, packageName types.PackageName, repositoryOwner string, repositoryName string, repositoryPackageRootPath string) (types.PackageName, error) {
	kurtosisYamlFilepath := path.Join(repositoryPackageRootPath, consts.DefaultKurtosisYamlFilename)
	repoGetContentOpts := &github.RepositoryContentGetOptions{
		Ref: "",
	}

	// get contents of kurtosis yaml file from GitHub
	kurtosisYamlFileContentResult, _, resp, err := packageExistRule.gitHubClient.Repositories.GetContents(ctx, repositoryOwner, repositoryName, kurtosisYamlFilepath, repoGetContentOpts)
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
