package files

import "github.com/rogueserenity/stenciler/config"

// CopyTemplated copies all templated files (those no located in raw-copy) and optionaly excluding files listed in
// init-only by passing them through the template engine, then writing the copies into the current working directory.
func CopyTemplated(repoDir string, template *config.Template, excludeInitOnly bool) error {
	return nil
}
