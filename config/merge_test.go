package config_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/rogueserenity/stenciler/config"
)

type MergeTestSuite struct {
	suite.Suite
}

func TestMergeTestSuite(t *testing.T) {
	suite.Run(t, new(MergeTestSuite))
}

func (s *MergeTestSuite) TestMergeRepoTemplateWithNoParams() {
	repo := &config.Template{
		Directory:           "foo",
		InitOnlyPaths:       []string{"init1"},
		RawCopyPaths:        []string{"raw1"},
		PreInitHookPaths:    []string{"pre-init1"},
		PostInitHookPaths:   []string{"post-init1"},
		PreUpdateHookPaths:  []string{"pre-update1"},
		PostUpdateHookPaths: []string{"post-update1"},
	}
	local := &config.Template{
		Repository: "https://github.com/owner/repo.git",
		Directory:  "foo",
		Params: []*config.Param{
			{
				Name:  "param1",
				Value: "value1",
			},
		},
		InitOnlyPaths:       []string{"init2"},
		RawCopyPaths:        []string{"raw2"},
		PreInitHookPaths:    []string{"pre-init2"},
		PostInitHookPaths:   []string{"post-init2"},
		PreUpdateHookPaths:  []string{"pre-update2"},
		PostUpdateHookPaths: []string{"post-update2"},
	}
	expected := &config.Template{
		Repository:          "https://github.com/owner/repo.git",
		Directory:           "foo",
		InitOnlyPaths:       []string{"init1"},
		RawCopyPaths:        []string{"raw1"},
		PreInitHookPaths:    []string{"pre-init1"},
		PostInitHookPaths:   []string{"post-init1"},
		PreUpdateHookPaths:  []string{"pre-update1"},
		PostUpdateHookPaths: []string{"post-update1"},
	}
	actual := config.Merge(repo, local)
	s.Require().Equal(expected, actual)
}

func (s *MergeTestSuite) TestMergeRepoTemplateWithOnlyNewParams() {
	repo := &config.Template{
		Directory: "foo",
		Params: []*config.Param{
			{
				Name:    "param2",
				Prompt:  "prompt2",
				Default: "mine",
			},
		},
	}
	local := &config.Template{
		Repository: "https://github.com/owner/repo.git",
		Directory:  "foo",
		Params: []*config.Param{
			{
				Name:  "param1",
				Value: "value1",
			},
		},
		InitOnlyPaths:       []string{"init2"},
		RawCopyPaths:        []string{"raw2"},
		PreInitHookPaths:    []string{"pre-init2"},
		PostInitHookPaths:   []string{"post-init2"},
		PreUpdateHookPaths:  []string{"pre-update2"},
		PostUpdateHookPaths: []string{"post-update2"},
	}
	expected := &config.Template{
		Repository: "https://github.com/owner/repo.git",
		Directory:  "foo",
		Params: []*config.Param{
			{
				Name:    "param2",
				Prompt:  "prompt2",
				Default: "mine",
			},
		},
	}
	actual := config.Merge(repo, local)
	s.Require().Equal(expected, actual)
}

func (s *MergeTestSuite) TestMergeTemplatesHaveMatchingParams() {
	repo := &config.Template{
		Directory: "foo",
		Params: []*config.Param{
			{
				Name:           "param1",
				Prompt:         "prompt2",
				Default:        "mine",
				ValidationHook: "hook1",
			},
		},
	}
	local := &config.Template{
		Repository: "https://github.com/owner/repo.git",
		Directory:  "foo",
		Params: []*config.Param{
			{
				Name:    "param1",
				Prompt:  "prompt1",
				Default: "yours",
				Value:   "value1",
			},
		},
	}
	expected := &config.Template{
		Repository: "https://github.com/owner/repo.git",
		Directory:  "foo",
		Params: []*config.Param{
			{
				Name:           "param1",
				Prompt:         "prompt2",
				Default:        "mine",
				ValidationHook: "hook1",
				Value:          "value1",
			},
		},
	}
	actual := config.Merge(repo, local)
	s.Require().Equal(expected, actual)
}

func (s *MergeTestSuite) TestMergeTemplatesHaveMatchingParamsLocalNoPrompt() {
	repo := &config.Template{
		Directory: "foo",
		Params: []*config.Param{
			{
				Name:           "param1",
				Prompt:         "prompt2",
				Default:        "mine",
				ValidationHook: "hook1",
			},
		},
	}
	local := &config.Template{
		Repository: "https://github.com/owner/repo.git",
		Directory:  "foo",
		Params: []*config.Param{
			{
				Name:  "param1",
				Value: "value1",
			},
		},
	}
	expected := &config.Template{
		Repository: "https://github.com/owner/repo.git",
		Directory:  "foo",
		Params: []*config.Param{
			{
				Name:           "param1",
				Prompt:         "prompt2",
				Default:        "mine",
				ValidationHook: "hook1",
			},
		},
	}
	actual := config.Merge(repo, local)
	s.Require().Equal(expected, actual)
}
