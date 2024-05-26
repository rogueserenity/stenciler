package cmd

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/rogueserenity/stenciler/config"
	"github.com/rogueserenity/stenciler/git"
)

// Command represents the init command
var initCmd = &cobra.Command{
	Args:  cobra.ExactArgs(1),
	Use:   "init [repoURL]",
	Short: "initialize a repository with the specified template",
	Long: `Initializes the current directory with the contents of the specified template.
Either the template repo URL should be specified as an argument or --template-directory should be specified.
The proper order of setting up a new repository is:
1. Create the repository with git
2. Intialize the repository with stenciler`,

	Run: func(cmd *cobra.Command, args []string) {
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
	if len(repoDir) == 0 {
		repoDir, err := git.Clone(repoURL.String(), authToken)
		if err != nil {
			cobra.CheckErr(err)
		}
		defer os.RemoveAll(repoDir)
		fmt.Printf("cloned %s to %s\n", repoURL, repoDir)
	}

	cfgFile := filepath.Join(repoDir, configFileName)
	cfg, err := config.ReadFromFile(cfgFile)
	if err != nil {
		cobra.CheckErr(err)
	}

	if len(cfg.Templates) == 0 {
		cobra.CheckErr("no templates found in config file")
	}

	localConfig := &config.Config{}

	if len(cfg.Templates) == 1 {
		localConfig.Templates = cfg.Templates
	} else {
		localConfig.Templates = append(localConfig.Templates, selectTemplate(cfg, cfgFile))
	}

	fmt.Println("init called with", repoURL)
}

func selectTemplate(cfg *config.Config, cfgFile string) config.Template {
	fmt.Printf("found %d templates in config file: %s\n", len(cfg.Templates), cfgFile)
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
		for _, t := range cfg.Templates {
			if t.Directory == d {
				template = &t
				break
			}
		}
	}
	return *template
}
