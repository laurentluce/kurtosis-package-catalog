package main

import (
	"github.com/kurtosis-tech/kurtosis-package-catalog/catalog-validator/importer"
	"github.com/kurtosis-tech/stacktrace"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
	"strings"
)

const (
	successExitCode = 0
	failureExitCode = 1

	forceColors   = true
	fullTimestamp = true

	logMethodAlongWithLogLine = true
	functionPathSeparator     = "."
	emptyFunctionName         = ""
)

func main() {

	configureLogger()

	packageCatalogYamlFilepath, err := getKurtosisPackageCatalogYAMLFilepathFromArgs()
	if err != nil {
		exitFailure(err)
	}

	packageCatalog, err := importer.ReadCatalog(packageCatalogYamlFilepath)
	if err != nil {
		exitFailure(err)
	}

	logrus.Infof("Package catalog is '%+v'", packageCatalog)

	logrus.Exit(successExitCode)
}

func getKurtosisPackageCatalogYAMLFilepathFromArgs() (string, error) {
	args := os.Args
	if len(args) < 2 {
		return "", stacktrace.NewError("expected to received the kurtosis package catalog YAML filepath as the first argument, but it was not received")
	}
	return args[1], nil
}

func configureLogger() {
	logrus.SetLevel(logrus.DebugLevel)
	// This allows the filename & function to be reported
	logrus.SetReportCaller(logMethodAlongWithLogLine)
	// NOTE: we'll want to change the ForceColors to false if we ever want structured logging
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:               forceColors,
		DisableColors:             false,
		ForceQuote:                false,
		DisableQuote:              false,
		EnvironmentOverrideColors: false,
		DisableTimestamp:          false,
		FullTimestamp:             fullTimestamp,
		TimestampFormat:           "",
		DisableSorting:            false,
		SortingFunc:               nil,
		DisableLevelTruncation:    false,
		PadLevelText:              false,
		QuoteEmptyFields:          false,
		FieldMap:                  nil,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			fullFunctionPath := strings.Split(f.Function, functionPathSeparator)
			functionName := fullFunctionPath[len(fullFunctionPath)-1]
			_, filename := path.Split(f.File)
			return emptyFunctionName, formatFilenameFunctionForLogs(filename, functionName)
		},
	})
}

func formatFilenameFunctionForLogs(filename string, functionName string) string {
	var output strings.Builder
	output.WriteString("[")
	output.WriteString(filename)
	output.WriteString(":")
	output.WriteString(functionName)
	output.WriteString("]")
	return output.String()
}

func exitFailure(err error) {
	logrus.Error(err.Error())
	logrus.Exit(failureExitCode)
}
