package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/rogueserenity/stenciler/config"
)

type ValidateTestSuite struct {
	suite.Suite

	templates []config.Template
}

func TestValidateTestSuite(t *testing.T) {
	suite.Run(t, new(ValidateTestSuite))
}

func (s *ValidateTestSuite) SetupTest() {
	s.templates = []config.Template{
		{
			Params: []*config.Param{
				{
					ValidationHook: "invalid-hook",
				},
			},
		},
		{
			PreInitHookPaths: []string{"invalid-hook"},
		},
		{
			PostInitHookPaths: []string{"invalid-hook"},
		},
		{
			PreUpdateHookPaths: []string{"invalid-hook"},
		},
		{
			PostUpdateHookPaths: []string{"invalid-hook"},
		},
	}
}

func (s *ValidateTestSuite) TestValidateWithNonExistentHooks() {
	for _, template := range s.templates {
		repoPath := "test-repo"
		err := template.Validate(repoPath)
		s.Require().ErrorContains(err, "hook invalid-hook does not exist")
	}
}

func (s *ValidateTestSuite) TestValidateWithNonExecutableHooks() {
	repoPath, err := os.MkdirTemp("", "test-repo")
	s.Require().NoError(err)
	defer os.RemoveAll(repoPath)
	file, err := os.Create(filepath.Join(repoPath, "invalid-hook"))
	s.Require().NoError(err)
	err = file.Close()
	s.Require().NoError(err)

	for _, template := range s.templates {
		err = template.Validate(repoPath)
		s.Require().ErrorContains(err, "hook invalid-hook is not executable")
	}
}

func (s *ValidateTestSuite) TestValidateWithValidHooks() {
	repoPath, err := os.MkdirTemp("", "test-repo")
	s.Require().NoError(err)
	defer os.RemoveAll(repoPath)
	file, err := os.Create(filepath.Join(repoPath, "invalid-hook"))
	s.Require().NoError(err)
	err = file.Chmod(0755)
	s.Require().NoError(err)
	err = file.Close()
	s.Require().NoError(err)

	for _, template := range s.templates {
		err = template.Validate(repoPath)
		s.Require().NoError(err)
	}
}
