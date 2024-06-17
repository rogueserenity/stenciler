package config_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/rogueserenity/stenciler/config"
)

var _ = DescribeTable("Merge", func(repoTemplate, localTemplate, expectedTemplate *config.Template) {
	mergedTemplate := config.Merge(repoTemplate, localTemplate)
	Expect(mergedTemplate).To(Equal(expectedTemplate))
},
	Entry("repo template has no params", &config.Template{
		Directory:           "foo",
		InitOnlyPaths:       []string{"init1"},
		RawCopyPaths:        []string{"raw1"},
		PreInitHookPaths:    []string{"pre-init1"},
		PostInitHookPaths:   []string{"post-init1"},
		PreUpdateHookPaths:  []string{"pre-update1"},
		PostUpdateHookPaths: []string{"post-update1"},
	}, &config.Template{
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
	}, &config.Template{
		Repository:          "https://github.com/owner/repo.git",
		Directory:           "foo",
		InitOnlyPaths:       []string{"init1"},
		RawCopyPaths:        []string{"raw1"},
		PreInitHookPaths:    []string{"pre-init1"},
		PostInitHookPaths:   []string{"post-init1"},
		PreUpdateHookPaths:  []string{"pre-update1"},
		PostUpdateHookPaths: []string{"post-update1"},
	}),

	Entry("repo template has only new params", &config.Template{
		Directory: "foo",
		Params: []*config.Param{
			{
				Name:    "param2",
				Prompt:  "prompt2",
				Default: "mine",
			},
		},
	}, &config.Template{
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
	}, &config.Template{
		Repository: "https://github.com/owner/repo.git",
		Directory:  "foo",
		Params: []*config.Param{
			{
				Name:    "param2",
				Prompt:  "prompt2",
				Default: "mine",
			},
		},
	}),

	Entry("templates have matching params", &config.Template{
		Directory: "foo",
		Params: []*config.Param{
			{
				Name:           "param1",
				Prompt:         "prompt2",
				Default:        "mine",
				ValidationHook: "hook1",
			},
		},
	}, &config.Template{
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
	}, &config.Template{
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
	}),
	Entry("templates have matching params but local does not have prompt", &config.Template{
		Directory: "foo",
		Params: []*config.Param{
			{
				Name:           "param1",
				Prompt:         "prompt2",
				Default:        "mine",
				ValidationHook: "hook1",
			},
		},
	}, &config.Template{
		Repository: "https://github.com/owner/repo.git",
		Directory:  "foo",
		Params: []*config.Param{
			{
				Name:  "param1",
				Value: "value1",
			},
		},
	}, &config.Template{
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
	}),
)
