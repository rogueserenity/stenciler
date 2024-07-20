package cmd

import (
	"fmt"
	"os"

	"github.com/carlmjohnson/versioninfo"
	"github.com/spf13/cobra"
)

const configFileName = ".stenciler.yaml"

// Persistent flags.
var (
	repoDir   string
	authToken string
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "stenciler",
	Short: "repository templates made easy",
	Long: `stenciler supports both initial templating of a repository and keeping
that repo up to date with changes from the repository`,

	PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
		if len(repoDir) > 0 {
			info, err := os.Stat(repoDir)
			if err != nil {
				return err
			}
			if !info.IsDir() {
				return fmt.Errorf("%s exists but is not a directory", repoDir)
			}
		}
		return nil
	},
}

// SetReleaseVersion sets the version of the application.
func SetReleaseVersion(version string) {
	if len(version) > 0 {
		rootCmd.Version = version
	} else {
		rootCmd.Version = versioninfo.Version
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&repoDir,
		"template-repo-dir",
		"r",
		"",
		"local directory to use as the template repository",
	)

	rootCmd.PersistentFlags().StringVarP(
		&authToken,
		"auth-token",
		"t",
		"",
		"authentication token for private remote repositories",
	)

	rootCmd.MarkFlagsMutuallyExclusive("template-repo-dir", "auth-token")
}
