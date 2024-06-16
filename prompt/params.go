package prompt

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rogueserenity/stenciler/config"
)

// ForParamValues prompts the user for values for any parameters that have a prompt defined and does not
// currently have a value associated. It will validate all values for params with a prompt defined.
func ForParamValues(template *config.Template, repoDir string) error {
	return ForParamValuesWithInOut(template, repoDir, os.Stdin, os.Stdout)
}

// ForParamValuesWithInOut prompts the user for values for any parameters that have a prompt defined and does not
// currently have a value associated. It will validate all values for params with a prompt defined.
func ForParamValuesWithInOut(template *config.Template, repoDir string, in io.Reader, out io.Writer) error {
	bufIn := bufio.NewReader(in)

	for _, p := range template.Params {
		if err := processParam(p, repoDir, bufIn, out); err != nil {
			return fmt.Errorf("error processing param: %w", err)
		}
	}

	return nil
}

func processParam(param *config.Param, repoDir string, in *bufio.Reader, out io.Writer) error {
	if len(param.Prompt) == 0 {
		return nil
	}

	if len(param.Value) == 0 {
		printParamPrompt(param, out)
		val, err := readParamPromptResponse(in)
		if err != nil {
			return fmt.Errorf("error reading response: %w", err)
		}
		if len(val) == 0 {
			val = param.Default
		}
		param.Value = val
	}

	return param.Validate(repoDir)
}

func printParamPrompt(param *config.Param, out io.Writer) {
	fmt.Fprint(out, param.Prompt)
	if len(param.Default) > 0 {
		fmt.Fprintf(out, " [%s]", param.Default)
	}
	fmt.Fprint(out, ": ")
}

func readParamPromptResponse(reader *bufio.Reader) (string, error) {
	value, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(value), nil
}
