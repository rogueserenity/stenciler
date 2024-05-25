package cmd

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

// Command represents the init command
var initCmd = &cobra.Command{
	Args:  cobra.MaximumNArgs(1),
	Use:   "init [repoURL]",
	Short: "initialize a repository with the specified template",
	Long: `Initializes the current directory with the contents of the specified template.
Either the template repo URL should be specified as an argument or --template-directory should be specified.
The proper order of setting up a new repository is:
1. Create the repository with git
2. Intialize the repository with stenciler`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			url, err := url.Parse(args[0])
			if err != nil {
				cobra.CheckErr("repoURL must be a valid URL")
			}

			fmt.Println("init called with", url)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}


