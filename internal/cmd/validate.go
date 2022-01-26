package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mt-sre/addon-metadata-operator/pkg/extractor"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
	"github.com/mt-sre/addon-metadata-operator/pkg/validate"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

func init() {
	validateCmd.Flags().StringVar(&validateEnv, "env", validateEnv, "integration, stage or production")
	validateCmd.Flags().StringVar(&validateVersion, "version", validateVersion, "addon imageset version")
	validateCmd.Flags().StringVar(&validateDisabled, "disabled", validateDisabled, "Disable specific validators, separated by ','. Can't be combined with --enabled.")
	validateCmd.Flags().StringVar(&validateEnabled, "enabled", validateEnabled, "Enable specific validators, separated by ','. Can't be combined with --disabled.")
	mtcli.AddCommand(validateCmd)
}

var (
	validateEnv      = "stage"
	validateVersion  = ""
	validateDisabled = ""
	validateEnabled  = ""
	validateExamples = []string{
		"  # Validate an addon in staging. Uses the latest version if it supports imageset.",
		"  mtcli validate --env stage --version latest internal/testdata/addons-imageset/reference-addon",
		"  # Validate a version 1.0.0 of a production addon using imageset.",
		"  mtcli validate --env production --version 1.0.0 <path/to/addon_dir>",
		"  # Validate a staging addon that is not using imageset, but a static indexImage.",
		"  mtcli validate --env stage <path/to/addon_dir>",
		"  # Validate an integration addon using imageset, disabling validators 001_foo and 002_bar.",
		"  mtcli validate --env integration --disabled AM0001,AM0002 <path/to/addon_dir>",
		"  # Validate an integration addon using imageset, enabled only 001_foo.",
		"  mtcli validate --env integration --enabled AM0001 <path/to/addon_dir>",
	}
	validateLong = "Validate an addon metadata and it's bundles against custom validators."
	validateCmd  = &cobra.Command{
		Use:     "validate",
		Short:   "Validate addon metadata, bundles and imagesets.",
		Long:    validateLong,
		Example: strings.Join(validateExamples, "\n"),
		Args:    cobra.ExactArgs(1),
		Run:     validateMain,
	}
)

func validateMain(cmd *cobra.Command, args []string) {
	addonDir, err := parseAddonDir(args[0])
	if err != nil {
		fail(1, "unable to parse the provided directory '%s': %v", args[0], err)
	}

	if err := verifyArgsAndFlags(addonDir); err != nil {
		fail(1, "unable to process flag or argument: %v", err)
	}

	meta, err := utils.NewMetaLoader(addonDir, validateEnv, validateVersion).Load()
	if err != nil {
		fail(1, "unable to load addon metadata from file '%v': %v", addonDir, err)
	}

	extractor := extractor.New()
	bundles, err := extractor.ExtractBundles(*meta.IndexImage, meta.OperatorName)
	if err != nil {
		fail(1, "unable to extract and parse bundles from the given index image: %v", err)
	}

	filter, err := validate.NewFilter(validateDisabled, validateEnabled)
	if err != nil {
		fail(1, "unable to process filter flags: %v", err)
	}

	success, errs := validate.ValidateCLI(*types.NewMetaBundle(meta, bundles), filter)
	if len(errs) > 0 {
		utils.PrintValidationErrors(errs)
		os.Exit(1)
	}

	if !success {
		os.Exit(1)
	}
}

func parseAddonDir(dir string) (string, error) {
	if !path.IsAbs(dir) {
		return filepath.Abs(dir)
	}
	return dir, nil
}

func fail(code int, msg string, args ...interface{}) {
	fmt.Printf("A fatal error occurred while preparing validations: "+msg+"\n", args...)

	os.Exit(code)
}

func verifyArgsAndFlags(addonDir string) error {
	if err := verifyAddonDir(addonDir); err != nil {
		return err
	}
	if err := verifyEnv(validateEnv); err != nil {
		return err
	}
	return verifyVersion(validateVersion)
}

// addonDir is an absolute path at this point
func verifyAddonDir(addonDir string) error {
	dir, err := os.Stat(addonDir)
	if err != nil {
		return fmt.Errorf("error while reading directory: %w", err)
	}
	if !dir.IsDir() {
		return fmt.Errorf("'%s' is not a directory", addonDir)
	}
	return nil
}

func verifyEnv(env string) error {
	if env != "integration" && env != "stage" && env != "production" {
		return fmt.Errorf("'%s' is not a valid environment; must be one of 'integration', 'stage' or 'production'", env)
	}
	return nil
}

func verifyVersion(version string) error {
	// unset version is OK, will fallback to meta.addonImageSetVersion
	if version == "" {
		return nil
	}
	// semver.IsValid(...) requires the following format vMAJOR.MINOR.PATCH
	// so we temporarily prefix the 'v' character
	if version != "latest" && !semver.IsValid(fmt.Sprintf("v%v", version)) {
		return fmt.Errorf("'%s' is not a valid version; must be one of 'latest' or match 'MAJOR.MINOR.PATCH'", version)
	}
	return nil
}
