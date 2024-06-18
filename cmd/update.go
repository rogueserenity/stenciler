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

// Command represents the init command.
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "updates a repository with the specified template",
	Long:  "Updates the current directory with the contents of the specified template.",

	Run: func(_ *cobra.Command, _ []string) {
		doUpdate()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func doUpdate() {
	slog.Debug("update called",
		slog.String("repoDir", repoDir),
		slog.Bool("authTokenProvided", len(authToken) > 0),
	)

	localTemplate := getLocalTemplateConfig()

	cloneRepo(localTemplate.Repository)

	repoTemplate := getRepoTemplateConfig(localTemplate.Directory)

	mergedTemplate := config.Merge(repoTemplate, localTemplate)

	err := mergedTemplate.Validate(repoDir)
	if err != nil {
		cobra.CheckErr(err)
	}

	err = prompt.ForParamValues(mergedTemplate, repoDir)
	if err != nil {
		cobra.CheckErr(err)
	}

	updateWrite(mergedTemplate)
}

func getLocalTemplateConfig() *config.Template {
	cfgFile := configFileName
	slog.Debug("config file path",
		slog.String("cfgFile", cfgFile),
		slog.String("repoDir", repoDir),
	)
	cfg, err := config.ReadFromFile(cfgFile)
	if err != nil {
		slog.Error("failed to read config file", slog.Any("error", err))
		cobra.CheckErr(err)
	}
	slog.Debug("config",
		slog.Any("config", cfg),
	)
	return cfg.Templates[0]
}

func getRepoTemplateConfig(templateDir string) *config.Template {
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
	slog.Debug("config",
		slog.Any("config", cfg),
	)

	template, err := prompt.SelectTemplate(templateDir, cfg)
	if err != nil {
		cobra.CheckErr(err)
	}

	return template
}

func cloneRepo(repoURL string) {
	if len(repoDir) == 0 {
		repoDir, err := git.Clone(repoURL, authToken)
		if err != nil {
			cobra.CheckErr(err)
		}
		defer os.RemoveAll(repoDir)
		slog.Debug("cloned repository",
			slog.String("repoDir", repoDir),
		)
	}
}

func updateWrite(template *config.Template) {
	localConfig := &config.Config{
		Templates: []*config.Template{template},
	}
	slog.Debug("writing config file", slog.Any("localConfig", localConfig))
	err := localConfig.WriteToFile(configFileName)
	if err != nil {
		cobra.CheckErr(err)
	}

	err = template.ExecuteHooks(repoDir, config.PreUpdateHook)
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

	err = template.ExecuteHooks(repoDir, config.PostUpdateHook)
	if err != nil {
		cobra.CheckErr(err)
	}
}
