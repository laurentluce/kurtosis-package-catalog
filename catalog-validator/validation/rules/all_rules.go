package rules

func GetAll() []Rule {
	allRules := []Rule{
		newDuplicatedPackageRule(),
	}

	return allRules
}
