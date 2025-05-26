package config_test

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/suite"

	"github.com/rogueserenity/stenciler/config"
)

type ConfigTestSuite struct {
	suite.Suite

	configText string
	cfg        *config.Config
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

func (s *ConfigTestSuite) SetupTest() {
	s.configText = `templates:
- repository: https://github.com/rogueserenity/stenciler-test
  directory: test
  params:
  - name: param1
    prompt: prompt1
    default: mine
    validation-hook: hook1
    value: yours
  init-only: ["init1"]
  raw-copy: ["raw1"]
  pre-init-hooks: ["pre-init1"]
  post-init-hooks: ["post-init1"]
  pre-update-hooks: ["pre-update1"]
  post-update-hooks: ["post-update1"]
`

	s.cfg = &config.Config{
		Templates: []*config.Template{
			{
				Repository: "https://github.com/rogueserenity/stenciler-test",
				Directory:  "test",
				Params: []*config.Param{
					{
						Name:           "param1",
						Prompt:         "prompt1",
						Default:        "mine",
						ValidationHook: "hook1",
						Value:          "yours",
					},
				},
				InitOnlyPaths:       []string{"init1"},
				RawCopyPaths:        []string{"raw1"},
				PreInitHookPaths:    []string{"pre-init1"},
				PostInitHookPaths:   []string{"post-init1"},
				PreUpdateHookPaths:  []string{"pre-update1"},
				PostUpdateHookPaths: []string{"post-update1"},
			},
		},
	}
}

func (s *ConfigTestSuite) TestRead() {
	reader := strings.NewReader(s.configText)
	actual, err := config.Read(reader)
	s.Require().NoError(err)
	s.Require().Equal(s.cfg, actual)
}

func (s *ConfigTestSuite) TestWrite() {
	writer := &strings.Builder{}
	err := s.cfg.Write(writer)
	s.Require().NoError(err)
	s.Require().YAMLEq(s.configText, writer.String())
}

func (s *ConfigTestSuite) TestParamValidateNoHook() {
	param := &config.Param{
		Name:    "test",
		Prompt:  "Test Prompt",
		Default: "default",
		Value:   "value",
	}

	err := param.Validate(faker.Word())
	s.Require().NoError(err)
}

func (s *ConfigTestSuite) TestParamValidateMissingHook() {
	param := &config.Param{
		Name:           "test",
		Prompt:         "Test Prompt",
		Default:        "default",
		Value:          "value",
		ValidationHook: "missing-hook",
	}

	dir := os.TempDir()
	err := param.Validate(dir)
	s.Require().ErrorContains(err, "failed to execute validation hook missing-hook on test with value value:")
}

func (s *ConfigTestSuite) TestParamValidateWithHook() {
	param := &config.Param{
		Name:           "test",
		Prompt:         "Test Prompt",
		Default:        "default",
		Value:          "value",
		ValidationHook: "valid.sh",
	}

	contents := `#!/bin/sh
	echo "validated_value"
`
	dir := os.TempDir()
	f, err := os.Create(path.Join(dir, "valid.sh"))
	s.Require().NoError(err)
	_, err = f.Write([]byte(contents))
	s.Require().NoError(err)
	s.Require().NoError(f.Close())
	defer os.Remove(path.Join(dir, "valid.sh"))

	err = param.Validate(dir)
	s.Require().NoError(err)

	s.Require().Equal("validated_value", param.Value)
}

func (s *ConfigTestSuite) TestExecuteHooksInvalidClass() {
	template := &config.Template{
		Repository: faker.URL(),
		Directory:  faker.Word(),
	}

	err := template.ExecuteHooks("test-dir", config.HookClass(42))
	s.Require().ErrorContains(err, "unknown hook class 42")
}
