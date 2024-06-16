package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/rogueserenity/stenciler/config"
	"github.com/rogueserenity/stenciler/files"
	"github.com/rogueserenity/stenciler/git"
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
	template := cfg.Templates[0]

	if len(repoDir) == 0 {
		repoDir, err = git.Clone(template.Repository, authToken)
		if err != nil {
			cobra.CheckErr(err)
		}
		defer os.RemoveAll(repoDir)
		slog.Debug("cloned repository",
			slog.String("repoDir", repoDir),
		)
	}

	err = template.Validate(repoDir)
	if err != nil {
		cobra.CheckErr(err)
	}

	updateWrite(template)
}

func updateWrite(template *config.Template) {
	err := template.ExecuteHooks(repoDir, config.PreUpdateHook)
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
