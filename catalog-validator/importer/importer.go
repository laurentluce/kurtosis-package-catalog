package importer

import (
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/catalog"
	"github.com/kurtosis-tech/stacktrace"
	"io"
	"net/http"
	"os"
)

const (
	currentPackageCatalogYamlFileURL = "https://raw.githubusercontent.com/kurtosis-tech/kurtosis-package-catalog/main/kurtosis-package-catalog.yml"
)

// GetNewPackageInTheCatalog compares the current state of the catalog in the repository main branch
// with the catalog read from a filepath and returns a subset containing the new packages being added
func GetNewPackageInTheCatalog(kurtosisPackageCatalogYamlFilepath string) (catalog.PackageCatalog, error) {
	newCatalog, err := readCatalog(kurtosisPackageCatalogYamlFilepath)
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred reading the catalog from '%s'", kurtosisPackageCatalogYamlFilepath)
	}

	currentCatalog, err := getCurrentPackageCatalog()
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred reading the current package catalog")
	}

	currentCatalogSet := map[string]bool{}

	for _, kurtosisPackage := range currentCatalog {
		kurtosisPackageStr := string(kurtosisPackage.GetPackageName())
		currentCatalogSet[kurtosisPackageStr] = true
	}

	var catalogWithNewPackages catalog.PackageCatalog

	for _, kurtosisPackage := range newCatalog {
		kurtosisPackageStr := string(kurtosisPackage.GetPackageName())
		if _, found := currentCatalogSet[kurtosisPackageStr]; !found {
			catalogWithNewPackages = append(catalogWithNewPackages, kurtosisPackage)
		}
	}
	return catalogWithNewPackages, nil
}

func readCatalog(kurtosisPackageCatalogYamlFilepath string) (catalog.PackageCatalog, error) {
	_, err := os.Stat(kurtosisPackageCatalogYamlFilepath)
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred checking for Kurtosis package catalog YAML file existence on '%s'", kurtosisPackageCatalogYamlFilepath)
	}

	fileBytes, err := os.ReadFile(kurtosisPackageCatalogYamlFilepath)
	if err != nil {
		return nil, stacktrace.Propagate(err, "attempted to read file with path '%v' but failed", kurtosisPackageCatalogYamlFilepath)
	}

	packageCatalog, err := catalog.GetPackageCatalogFromYamlFileContent(fileBytes)
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred reading the Kurtosis package catalog YAML file content from '%s'", kurtosisPackageCatalogYamlFilepath)
	}

	return packageCatalog, nil
}

func getCurrentPackageCatalog() (catalog.PackageCatalog, error) {

	response, getErr := http.Get(currentPackageCatalogYamlFileURL)
	if getErr != nil {
		return nil, stacktrace.Propagate(getErr, "an error occurred getting the yaml file content from URL '%s'", currentPackageCatalogYamlFileURL)
	}
	defer response.Body.Close()
	responseBodyBytes, readAllErr := io.ReadAll(response.Body)
	if readAllErr != nil {
		return nil, stacktrace.Propagate(readAllErr, "an error occurred reading the yaml file content")
	}

	packageCatalog, err := catalog.GetPackageCatalogFromYamlFileContent(responseBodyBytes)
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred reading the Kurtosis package catalog YAML file content from '%s'", currentPackageCatalogYamlFileURL)
	}

	return packageCatalog, nil
}
