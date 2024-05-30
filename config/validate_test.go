package config_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/rogueserenity/stenciler/config"
)

var _ = Describe("Validate", func() {

	DescribeTable("Validate with non-existent hooks", func(template config.Template) {
		repoPath := "test-repo"
		err := template.Validate(repoPath)
		Expect(err).To(MatchError("hook invalid-hook does not exist"))
	},
		Entry("ValidationHook", config.Template{
			Params: []config.Param{
				{
					ValidationHook: "invalid-hook",
				},
			},
		}),
		Entry("PreInitHookPaths", config.Template{
			PreInitHookPaths: []string{"invalid-hook"},
		}),
		Entry("PostInitHookPaths", config.Template{
			PostInitHookPaths: []string{"invalid-hook"},
		}),
		Entry("PreUpdateHookPaths", config.Template{
			PreUpdateHookPaths: []string{"invalid-hook"},
		}),
		Entry("PostUpdateHookPaths", config.Template{
			PostUpdateHookPaths: []string{"invalid-hook"},
		}),
	)

	DescribeTable("Validate with non-executable hooks", func(template config.Template) {
		repoPath, err := os.MkdirTemp("", "test-repo")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(repoPath)
		file, err := os.Create(filepath.Join(repoPath, "invalid-hook"))
		Expect(err).NotTo(HaveOccurred())
		err = file.Close()
		Expect(err).NotTo(HaveOccurred())

		err = template.Validate(repoPath)
		Expect(err).To(MatchError("hook invalid-hook is not executable"))
	},
		Entry("ValidationHook", config.Template{
			Params: []config.Param{
				{
					ValidationHook: "invalid-hook",
				},
			},
		}),
		Entry("PreInitHookPaths", config.Template{
			PreInitHookPaths: []string{"invalid-hook"},
		}),
		Entry("PostInitHookPaths", config.Template{
			PostInitHookPaths: []string{"invalid-hook"},
		}),
		Entry("PreUpdateHookPaths", config.Template{
			PreUpdateHookPaths: []string{"invalid-hook"},
		}),
		Entry("PostUpdateHookPaths", config.Template{
			PostUpdateHookPaths: []string{"invalid-hook"},
		}),
	)

	DescribeTable("Validate with valid hooks", func(template config.Template) {
		repoPath, err := os.MkdirTemp("", "test-repo")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(repoPath)
		file, err := os.Create(filepath.Join(repoPath, "invalid-hook"))
		Expect(err).NotTo(HaveOccurred())
		err = file.Chmod(0755)
		Expect(err).NotTo(HaveOccurred())
		err = file.Close()
		Expect(err).NotTo(HaveOccurred())

		err = template.Validate(repoPath)
		Expect(err).ToNot(HaveOccurred())
	},
		Entry("ValidationHook", config.Template{
			Params: []config.Param{
				{
					ValidationHook: "invalid-hook",
				},
			},
		}),
		Entry("PreInitHookPaths", config.Template{
			PreInitHookPaths: []string{"invalid-hook"},
		}),
		Entry("PostInitHookPaths", config.Template{
			PostInitHookPaths: []string{"invalid-hook"},
		}),
		Entry("PreUpdateHookPaths", config.Template{
			PreUpdateHookPaths: []string{"invalid-hook"},
		}),
		Entry("PostUpdateHookPaths", config.Template{
			PostUpdateHookPaths: []string{"invalid-hook"},
		}),
	)

})
