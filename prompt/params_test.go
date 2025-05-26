package prompt_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/suite"

	"github.com/rogueserenity/stenciler/config"
	"github.com/rogueserenity/stenciler/prompt"
)

type PromptParamsTestSuite struct {
	suite.Suite

	repoDir string

	stdin  *bytes.Buffer
	stdout *strings.Builder
}

func TestPromptParamsTestSuite(t *testing.T) {
	suite.Run(t, new(PromptParamsTestSuite))
}

func (s *PromptParamsTestSuite) SetupTest() {
	s.repoDir = faker.Word()

	s.stdin = &bytes.Buffer{}
	s.stdout = &strings.Builder{}
}

func (s *PromptParamsTestSuite) TestForParamValuesNoPrompt() {
	template := &config.Template{
		Params: []*config.Param{
			{
				Name:  "foo",
				Value: "bar",
			},
		},
	}

	err := prompt.ForParamValuesWithInOut(template, s.repoDir, s.stdin, s.stdout)
	s.NoError(err)
	s.Empty(s.stdout.String())
}

func (s *PromptParamsTestSuite) TestForParamValuesWithPromptAndValue() {
	template := &config.Template{
		Params: []*config.Param{
			{
				Name:   "foo",
				Prompt: "Enter a value for foo",
				Value:  "bar",
			},
		},
	}

	err := prompt.ForParamValuesWithInOut(template, s.repoDir, s.stdin, s.stdout)
	s.NoError(err)
	s.Empty(s.stdout.String())
}

func (s *PromptParamsTestSuite) TestForParamValuesWithPromptAndNoValue() {
	template := &config.Template{
		Params: []*config.Param{
			{
				Name:   "foo",
				Prompt: "Enter a value for foo",
			},
		},
	}

	s.stdin.WriteString("baz\n")

	err := prompt.ForParamValuesWithInOut(template, s.repoDir, s.stdin, s.stdout)
	s.NoError(err)

	expectedOutput := "Enter a value for foo: "
	s.Equal(expectedOutput, s.stdout.String())
	s.Equal("baz", template.Params[0].Value)
}

func (s *PromptParamsTestSuite) TestForParamValuesWithPromptAndDefaultValue() {
	template := &config.Template{
		Params: []*config.Param{
			{
				Name:    "foo",
				Prompt:  "Enter a value for foo",
				Default: "default_value",
			},
		},
	}

	s.stdin.WriteString("\n")

	err := prompt.ForParamValuesWithInOut(template, s.repoDir, s.stdin, s.stdout)
	s.NoError(err)

	expectedOutput := "Enter a value for foo [default_value]: "
	s.Equal(expectedOutput, s.stdout.String())
	s.Equal("default_value", template.Params[0].Value)
}

func (s *PromptParamsTestSuite) TestForParamValuesWithPromptAndDefaultValueAndInput() {
	template := &config.Template{
		Params: []*config.Param{
			{
				Name:    "foo",
				Prompt:  "Enter a value for foo",
				Default: "default_value",
			},
		},
	}

	s.stdin.WriteString("input_value\n")

	err := prompt.ForParamValuesWithInOut(template, s.repoDir, s.stdin, s.stdout)
	s.NoError(err)

	expectedOutput := "Enter a value for foo [default_value]: "
	s.Equal(expectedOutput, s.stdout.String())
	s.Equal("input_value", template.Params[0].Value)
}
