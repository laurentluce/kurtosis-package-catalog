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
	"image"
	_ "image/png" // need to import it to get the PNG Encoder/Decoder
	"net/http"
	"path"
	"strings"
)

const (
	validPackageIconRuleName = "Valid package icon"
	minImageSize             = 120
	maxImageSize             = 1024
)

// validPackageIconRule checks if the package icon is valid by checking if:
// 1- if the png image exist, does not return an error if it not because it's not mandatory yet
// 2- if the image size is equal or bigger that the minImageSize
// 3- if the image size is equal or greater than maxImageSize
// 4- if the aspect ratio is 1:1 (a square image)
type validPackageIconRule struct {
	name         string
	gitHubClient *github.Client
}

func newValidPackageIconRule(gitHubClient *github.Client) *validPackageIconRule {
	return &validPackageIconRule{name: validPackageIconRuleName, gitHubClient: gitHubClient}
}

func (validPackageIconRule *validPackageIconRule) GetName() RuleName {
	return RuleName(validPackageIconRule.name)
}

func (validPackageIconRule *validPackageIconRule) Check(ctx context.Context, catalog catalog.PackageCatalog) *CheckResult {

	wasValidated := true
	failures := map[types.PackageName][]string{}

	for _, packageData := range catalog {
		packageName := packageData.GetPackageName()
		logrus.Debugf("Checking if package '%s' contains a valid icon...", packageName)
		repositoryOwner := packageData.GetRepositoryOwner()
		repositoryName := packageData.GetRepositoryName()
		repositoryPackageRootPath := packageData.GetRepositoryPackageRootPath()
		packageFailures := []string{}
		packageIconImageConfig, err := validPackageIconRule.getPackageIconImageConfig(ctx, packageName, repositoryOwner, repositoryName, repositoryPackageRootPath)
		if err != nil {
			errorFailure := fmt.Sprintf("an error occurred getting the Kurtosis package icon image config for package '%s'. Error was:\n%s", packageName, err.Error())
			packageFailures = append(packageFailures, errorFailure)
		}
		if err == nil && packageIconImageConfig == nil {
			logrus.Debugf("package '%s' does not have an icon yet.", packageName)
			continue
		}
		if err == nil {
			packageIconWidth := packageIconImageConfig.Width
			packageIconHeight := packageIconImageConfig.Height

			if packageIconWidth < minImageSize || packageIconHeight < minImageSize {
				invalidMinSizeMsg := fmt.Sprintf(
					"invalid image min size, it is smaller than expected. "+
						"Valid min value is '%dpx' and the current size is width: %dpx and height: %dpx",
					minImageSize,
					packageIconWidth,
					packageIconHeight,
				)
				packageFailures = append(packageFailures, invalidMinSizeMsg)
			}

			if packageIconWidth > maxImageSize || packageIconHeight > maxImageSize {
				invalidMaxSizeMsg := fmt.Sprintf(
					"invalid image max size, it is bigger than expected. "+
						"Valid max value is '%dpx' and the current size is width: %dpx and height: %dpx",
					maxImageSize,
					packageIconWidth,
					packageIconHeight,
				)
				packageFailures = append(packageFailures, invalidMaxSizeMsg)
			}

			if packageIconWidth != packageIconHeight {
				invalidAspectRatioMsg := "invalid aspect ratio, the accepted aspect ration is 1:1 (a square image)."

				packageFailures = append(packageFailures, invalidAspectRatioMsg)
			}
		}

		if len(packageFailures) > 0 {
			failures[packageName] = packageFailures
			wasValidated = false
			continue
		}
		logrus.Debugf("...package icon for '%s' successfully validated.", packageName)
	}

	checkResult := newCheckResult(validPackageIconRule.GetName(), wasValidated, failures)

	return checkResult
}

func (validPackageIconRule *validPackageIconRule) getPackageIconImageConfig(ctx context.Context, packageName types.PackageName, repositoryOwner string, repositoryName string, repositoryPackageRootPath string) (*image.Config, error) {
	packageIconFilepath := path.Join(repositoryPackageRootPath, consts.KurtosisPackageIconImgName)
	repoGetContentOpts := &github.RepositoryContentGetOptions{
		Ref: "",
	}

	// get contents of kurtosis package icon file from GitHub
	packageIconFileContentResult, _, resp, err := validPackageIconRule.gitHubClient.Repositories.GetContents(ctx, repositoryOwner, repositoryName, packageIconFilepath, repoGetContentOpts)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			// having the icon is not mandatory
			return nil, nil
		}
		return nil, stacktrace.Propagate(err, "an error occurred reading content of Kurtosis Package '%s' - file '%s'", packageName, packageIconFilepath)
	}

	rawPackageIconContentStr, err := packageIconFileContentResult.GetContent()
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred getting the '%s' base 64 file content in package '%s'", packageIconFilepath, packageName)
	}

	packageIconContentReader := strings.NewReader(rawPackageIconContentStr)

	packageIconConfig, _, err := image.DecodeConfig(packageIconContentReader)
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred while decoding the '%s' image file in package '%s'", packageIconFilepath, packageName)
	}

	return &packageIconConfig, nil
}
