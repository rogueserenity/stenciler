package prompt_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/rogueserenity/stenciler/config"
	"github.com/rogueserenity/stenciler/prompt"
)

type PromptForTemplateTestSuite struct {
	suite.Suite

	stdin  *bytes.Buffer
	stdout *strings.Builder
}

func TestPromptForTemplateTestSuite(t *testing.T) {
	suite.Run(t, new(PromptForTemplateTestSuite))
}

func (s *PromptForTemplateTestSuite) SetupTest() {
	s.stdin = &bytes.Buffer{}
	s.stdout = &strings.Builder{}
}

func (s *PromptForTemplateTestSuite) TestSelectTemplateWithValidDir() {
	templateDir := "foo"
	cfg := &config.Config{
		Templates: []*config.Template{
			{Directory: "foo"},
			{Directory: "bar"},
		},
	}

	template, err := prompt.SelectTemplateWithInOut(templateDir, cfg, s.stdin, s.stdout)
	s.NoError(err)
	s.NotNil(template)
	s.Equal("foo", template.Directory)
}

func (s *PromptForTemplateTestSuite) TestSelectTemplateWithInvalidDir() {
	templateDir := "baz"
	cfg := &config.Config{
		Templates: []*config.Template{
			{Directory: "foo"},
			{Directory: "bar"},
		},
	}

	template, err := prompt.SelectTemplateWithInOut(templateDir, cfg, s.stdin, s.stdout)
	s.ErrorContains(err, "template directory not found in config")
	s.Nil(template)
}

func (s *PromptForTemplateTestSuite) TestSelectTemplateWithSingleTemplate() {
	templateDir := ""
	cfg := &config.Config{
		Templates: []*config.Template{
			{Directory: "foo"},
		},
	}

	template, err := prompt.SelectTemplateWithInOut(templateDir, cfg, s.stdin, s.stdout)
	s.NoError(err)
	s.NotNil(template)
	s.Equal("foo", template.Directory)
}

func (s *PromptForTemplateTestSuite) TestSelectTemplateWithMultipleTemplatesInvalidSelection() {
	templateDir := ""
	cfg := &config.Config{
		Templates: []*config.Template{
			{Directory: "foo"},
			{Directory: "bar"},
		},
	}

	s.stdin.WriteString("baz\nfoo\n") // Invalid first selection, valid second selection

	template, err := prompt.SelectTemplateWithInOut(templateDir, cfg, s.stdin, s.stdout)
	s.NoError(err)
	s.NotNil(template)
	s.Equal("foo", template.Directory)

	expectedOutput := "Available templates:\n>  bar\n>  foo\nplease specify the template directory to use: " +
		"Available templates:\n>  bar\n>  foo\nplease specify the template directory to use: "
	s.Equal(expectedOutput, s.stdout.String())
}

func (s *PromptForTemplateTestSuite) TestSelectTemplateWithMultipleTemplatesValidSelection() {
	templateDir := ""
	cfg := &config.Config{
		Templates: []*config.Template{
			{Directory: "foo"},
			{Directory: "bar"},
		},
	}

	s.stdin.WriteString("foo\n") // Valid selection

	template, err := prompt.SelectTemplateWithInOut(templateDir, cfg, s.stdin, s.stdout)
	s.NoError(err)
	s.NotNil(template)
	s.Equal("foo", template.Directory)

	expectedOutput := "Available templates:\n>  bar\n>  foo\nplease specify the template directory to use: "
	s.Equal(expectedOutput, s.stdout.String())
}
