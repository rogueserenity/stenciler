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
	"github.com/rogueserenity/stenciler/hooks"
)

// Command represents the init command.
var initCmd = &cobra.Command{
	Args:  cobra.ExactArgs(1),
	Use:   "init [repoURL]",
	Short: "initialize a repository with the specified template",
	Long: `Initializes the current directory with the contents of the specified template.
Either the template repo URL should be specified as an argument or --template-directory should be specified.
The proper order of setting up a new repository is:
1. Create the repository with git
2. Intialize the repository with stenciler`,

	Run: func(_ *cobra.Command, args []string) {
		url, err := url.Parse(args[0])
		if err != nil {
			cobra.CheckErr("repoURL must be a valid URL")
		}
		doInit(url)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func doInit(repoURL *url.URL) {
	fmt.Println("init called with", repoURL)

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

	localConfig := &config.Config{}

	if len(cfg.Templates) == 1 {
		localConfig.Templates = cfg.Templates
	} else {
		localConfig.Templates = append(localConfig.Templates, selectTemplate(cfg, cfgFile))
	}
	template := &localConfig.Templates[0]

	err = hooks.Validate(template, repoDir)
	if err != nil {
		cobra.CheckErr(err)
	}

	prompt(template)

	initialWrite(localConfig)
}

func selectTemplate(cfg *config.Config, cfgFile string) config.Template {
	fmt.Printf("found %d templates in config file: %s\n", len(cfg.Templates), cfgFile)

	var templateMap = make(map[string]config.Template)
	for _, t := range cfg.Templates {
		templateMap[t.Directory] = t
	}

	var template *config.Template
	for template == nil {
		for _, t := range cfg.Templates {
			fmt.Println("directory:", t.Directory)
		}
		fmt.Print("please specify the template to use: ")
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

func prompt(template *config.Template) {
	reader := bufio.NewReader(os.Stdin)

	for _, p := range template.Params {
		if len(p.Prompt) == 0 {
			continue
		}
		fmt.Print(p.Prompt)
		if len(p.Default) > 0 {
			fmt.Printf(" [%s]", p.Default)
		}
		fmt.Print(": ")
		value, err := reader.ReadString('\n')
		if err != nil {
			cobra.CheckErr(err)
		}
		value = strings.TrimSpace(value)
		if len(value) == 0 {
			value = p.Default
		}

		if len(p.ValidationHook) > 0 {
			value, err = hooks.ExecuteValidationHook(filepath.Join(repoDir, p.ValidationHook), p.Name, value)
			if err != nil {
				cobra.CheckErr(err)
			}
		}

		p.Value = value
	}
}

func initialWrite(localConfig *config.Config) {
	err := localConfig.WriteToFile(configFileName)
	if err != nil {
		cobra.CheckErr(err)
	}

	template := &localConfig.Templates[0]

	err = hooks.ExecuteHooks(repoDir, template.PreInitHookPaths)
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

	err = hooks.ExecuteHooks(repoDir, template.PostInitHookPaths)
	if err != nil {
		cobra.CheckErr(err)
	}
}
