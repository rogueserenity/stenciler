package cmd

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/rogueserenity/stenciler/config"
	"github.com/rogueserenity/stenciler/files"
	"github.com/rogueserenity/stenciler/git"
	"github.com/rogueserenity/stenciler/prompt"
)

var (
	templateDir string
)

// Command represents the init command.
var initCmd = &cobra.Command{
	Args:  cobra.ExactArgs(1),
	Use:   "init repoURL",
	Short: "initialize a repository with the specified template",
	Long:  "Initializes the current directory with the contents of the specified template.",

	Run: func(_ *cobra.Command, args []string) {
		doInit(args[0])
	},
}

func init() {
	initCmd.Flags().StringVarP(
		&templateDir,
		"template-dir",
		"d",
		"",
		"template directory to use from the config file",
	)
	rootCmd.AddCommand(initCmd)
}

func doInit(repoURL string) {
	slog.Debug("init called",
		slog.String("repoURL", repoURL),
		slog.String("repoDir", repoDir),
		slog.Bool("authTokenProvided", len(authToken) > 0),
		slog.String("templateDir", templateDir),
	)

	var err error
	if len(repoDir) == 0 {
		repoDir, err = git.Clone(repoURL, authToken)
		if err != nil {
			cobra.CheckErr(err)
		}
		defer os.RemoveAll(repoDir)
		slog.Debug("cloned repository",
			slog.String("repoURL", repoURL),
			slog.String("repoDir", repoDir),
		)
	}

	cfgFile := filepath.Join(repoDir, configFileName)
	slog.Debug("config file path",
		slog.String("cfgFile", cfgFile),
		slog.String("repoDir", repoDir),
		slog.String("configFileName", configFileName),
	)
	cfg, err := config.ReadFromFile(cfgFile)
	if err != nil {
		cobra.CheckErr(err)
	}

	slog.Debug("config", slog.Any("config", *cfg))

	if len(cfg.Templates) == 0 {
		cobra.CheckErr("no templates found in config file")
	}

	template, err := prompt.SelectTemplate(templateDir, cfg)
	if err != nil {
		cobra.CheckErr(err)
	}
	localConfig := &config.Config{
		Templates: []*config.Template{template},
	}
	template.Repository = repoURL

	err = template.Validate(repoDir)
	if err != nil {
		cobra.CheckErr(err)
	}

	err = prompt.ForParamValues(template, repoDir)
	if err != nil {
		cobra.CheckErr(err)
	}

	initialWrite(localConfig)
}

func initialWrite(localConfig *config.Config) {
	slog.Debug("writing config file", slog.Any("localConfig", localConfig))
	err := localConfig.WriteToFile(configFileName)
	if err != nil {
		cobra.CheckErr(err)
	}

	template := localConfig.Templates[0]

	err = template.ExecuteHooks(repoDir, config.PreInitHook)
	if err != nil {
		cobra.CheckErr(err)
	}

	err = files.CopyRaw(repoDir, template)
	if err != nil {
		cobra.CheckErr(err)
	}

	err = files.CopyTemplated(repoDir, template)
	if err != nil {
		cobra.CheckErr(err)
	}

	err = template.ExecuteHooks(repoDir, config.PostInitHook)
	if err != nil {
		cobra.CheckErr(err)
	}
}
