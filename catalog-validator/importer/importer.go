package importer

import (
	"os"

	"github.com/kurtosis-tech/kurtosis-package-indexer/server/catalog"
	"github.com/kurtosis-tech/stacktrace"
)

const (
	kurtosisPackageCatalogYamlFilepath = "../kurtosis-package-catalog.yml"
)

func ReadCatalog() (catalog.PackageCatalog, error) {
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
