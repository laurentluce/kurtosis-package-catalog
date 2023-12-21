package rules

import "github.com/kurtosis-tech/kurtosis-package-indexer/server/catalog"

type Rule interface {
	GetName() string
	Check(catalog.PackageCatalog) error
}
