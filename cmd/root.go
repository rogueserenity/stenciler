package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/rogueserenity/stenciler/config"
)

const defaultConfigFile = ".stenciler.yaml"

// persistent flags
var (
	configFile string
	repoDir    string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "stenciler",
	Short: "repository templates made easy",
	Long: `stenciler supports both initial templating of a repository and keeping
that repo up to date with changes from the repository`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if len(repoDir) > 0 {
			fi, err := os.Stat(repoDir)
			if err != nil {
				return err
			}
			if !fi.IsDir() {
				return fmt.Errorf("%s exists but is not a directory", repoDir)
			}
		}
		return nil
	},
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
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "alternate config file (default is .stenciler.yaml)")
	rootCmd.PersistentFlags().StringVarP(&repoDir, "template-repo-dir", "t", "", "local template repository directory")
}

func loadConfig() (*config.Config, error) {
	cfgFile := defaultConfigFile
	if len(configFile) > 0 {
		cfgFile = configFile
	}

	return config.ReadFromFile(cfgFile)
}
