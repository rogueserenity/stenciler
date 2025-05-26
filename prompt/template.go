package prompt

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/rogueserenity/stenciler/config"
)

// SelectTemplate prompts the user to select a template if more than one template is defined in the configuration
// and no template directory is specified.
func SelectTemplate(templateDir string, cfg *config.Config) (*config.Template, error) {
	return SelectTemplateWithInOut(templateDir, cfg, os.Stdin, os.Stdout)
}

// SelectTemplateWithInOut prompts the user to select a template if more than one template is defined in the
// configuration and no template directory is specified. It uses the provided input and output streams.
func SelectTemplateWithInOut(
	templateDir string,
	cfg *config.Config,
	in io.Reader,
	out io.Writer) (*config.Template, error) {
	templateMap := createTemplateMap(cfg)

	if len(templateDir) > 0 {
		if t, ok := templateMap[templateDir]; ok {
			return t, nil
		}
		return nil, errors.New("template directory not found in config")
	}

	if len(cfg.Templates) == 1 {
		return cfg.Templates[0], nil
	}

	return promptForTemplateDir(templateMap, in, out)
}

func createTemplateMap(cfg *config.Config) map[string]*config.Template {
	var templateMap = make(map[string]*config.Template)
	for _, t := range cfg.Templates {
		templateMap[t.Directory] = t
	}
	return templateMap
}

func sortedKeys(m map[string]*config.Template) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

func promptForTemplateDir(
	templateMap map[string]*config.Template,
	in io.Reader,
	out io.Writer) (*config.Template, error) {
	reader := bufio.NewReader(in)

	keys := sortedKeys(templateMap)
	var template *config.Template
	for template == nil {
		fmt.Fprintln(out, "Available templates:")
		for _, k := range keys {
			fmt.Fprintln(out, "> ", k)
		}
		fmt.Fprint(out, "please specify the template directory to use: ")
		d, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed to read input: %w", err)
		}
		d = strings.TrimSpace(d)
		if t, ok := templateMap[d]; ok {
			template = t
		}
	}
	return template, nil
}
