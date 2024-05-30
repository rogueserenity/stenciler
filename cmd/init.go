package cmd

import (
	"bufio"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/rogueserenity/stenciler/config"
	"github.com/rogueserenity/stenciler/files"
	"github.com/rogueserenity/stenciler/git"
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
		url, err := url.Parse(args[0])
		if err != nil {
			cobra.CheckErr("repoURL must be a valid URL")
		}
		doInit(url)
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

func doInit(repoURL *url.URL) {
	slog.Debug("init called",
		slog.String("repoURL", repoURL.String()),
		slog.String("repoDir", repoDir),
		slog.Bool("authTokenProvided", len(authToken) > 0),
		slog.String("templateDir", templateDir),
	)

	var err error
	if len(repoDir) == 0 {
		repoDir, err = git.Clone(repoURL.String(), authToken)
		if err != nil {
			cobra.CheckErr(err)
		}
		defer os.RemoveAll(repoDir)
		slog.Debug("cloned repository",
			slog.String("repoURL", repoURL.String()),
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

	slog.Debug("config",
		slog.Any("config", cfg),
	)

	if len(cfg.Templates) == 0 {
		cobra.CheckErr("no templates found in config file")
	}

	localConfig := &config.Config{
		Templates: []config.Template{selectTemplate(cfg)},
	}
	template := &localConfig.Templates[0]
	template.Repository = repoURL.String()

	err = template.Validate(repoDir)
	if err != nil {
		cobra.CheckErr(err)
	}

	prompt(template)

	initialWrite(localConfig)
}

func selectTemplate(cfg *config.Config) config.Template {
	var templateMap = make(map[string]config.Template)
	for _, t := range cfg.Templates {
		templateMap[t.Directory] = t
	}

	if len(templateDir) > 0 {
		if t, ok := templateMap[templateDir]; ok {
			return t
		}
		cobra.CheckErr("template name not found in config")
	}

	if len(cfg.Templates) == 1 {
		return cfg.Templates[0]
	}

	var template *config.Template
	for template == nil {
		fmt.Println("Available templates:")
		for _, t := range cfg.Templates {
			fmt.Println("> ", t.Directory)
		}
		fmt.Print("please specify the template directory to use: ")
		reader := bufio.NewReader(os.Stdin)
		d, err := reader.ReadString('\n')
		if err != nil {
			cobra.CheckErr(err)
		}
		d = strings.TrimSpace(d)
		if t, ok := templateMap[d]; ok {
			template = &t
		}
	}
	return *template
}

func printPrompt(param config.Param) {
	fmt.Print(param.Prompt)
	if len(param.Default) > 0 {
		fmt.Printf(" [%s]", param.Default)
	}
	fmt.Print(": ")
}

func readPromptResponse() string {
	reader := bufio.NewReader(os.Stdin)
	value, err := reader.ReadString('\n')
	if err != nil {
		cobra.CheckErr(err)
	}
	return strings.TrimSpace(value)
}

func prompt(template *config.Template) {
	for _, p := range template.Params {
		if len(p.Prompt) == 0 {
			continue
		}

		printPrompt(p)
		p.Value = readPromptResponse()
		if len(p.Value) == 0 {
			p.Value = p.Default
		}

		p.Validate(repoDir)
	}
}

func initialWrite(localConfig *config.Config) {
	slog.Debug("writing config file", slog.Any("localConfig", localConfig))
	err := localConfig.WriteToFile(configFileName)
	if err != nil {
		cobra.CheckErr(err)
	}

	template := &localConfig.Templates[0]

	err = template.ExecuteHooks(repoDir, config.PreInitHook)
	if err != nil {
		cobra.CheckErr(err)
	}

	err = files.CopyRaw(repoDir, template)
	if err != nil {
		cobra.CheckErr(err)
	}

	err = files.CopyTemplated(repoDir, template, false)
	if err != nil {
		cobra.CheckErr(err)
	}

	err = template.ExecuteHooks(repoDir, config.PostInitHook)
	if err != nil {
		cobra.CheckErr(err)
	}
}
