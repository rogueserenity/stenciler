package config_test

import (
	"io"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/rogueserenity/stenciler/config"
)

var _ = Describe("Config", func() {

	Describe("Read", func() {
		var (
			reader io.Reader
			cfg    *config.Config
			err    error
		)

		JustBeforeEach(func() {
			cfg, err = config.Read(reader)
		})

		Context("when given a config with all possible fields", func() {
			var input = `
templates:
-
  repository: https://github.com/rogueserenity/stenciler-test
  directory: test
  params:
  -
    name: param1
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
			BeforeEach(func() {
				reader = strings.NewReader(input)
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return a Config object with all fields", func() {
				Expect(cfg).ToNot(BeNil())
				Expect(cfg.Templates).To(HaveLen(1))
				template := cfg.Templates[0]
				Expect(template.Repository).To(Equal("https://github.com/rogueserenity/stenciler-test"))
				Expect(template.Directory).To(Equal("test"))
				Expect(template.Params).To(HaveLen(1))
				param := template.Params[0]
				Expect(param.Name).To(Equal("param1"))
				Expect(param.Prompt).To(Equal("prompt1"))
				Expect(param.Default).To(Equal("mine"))
				Expect(param.ValidationHook).To(Equal("hook1"))
				Expect(param.Value).To(Equal("yours"))
				Expect(template.InitOnlyPaths).To(Equal([]string{"init1"}))
				Expect(template.RawCopyPaths).To(Equal([]string{"raw1"}))
				Expect(template.PreInitHookPaths).To(Equal([]string{"pre-init1"}))
				Expect(template.PostInitHookPaths).To(Equal([]string{"post-init1"}))
				Expect(template.PreUpdateHookPaths).To(Equal([]string{"pre-update1"}))
				Expect(template.PostUpdateHookPaths).To(Equal([]string{"post-update1"}))
			})
		})
	})

	Describe("Write", func() {
		var (
			cfg *config.Config

			writer *strings.Builder
			err    error

			expected string
			actual   string
		)

		BeforeEach(func() {
			writer = &strings.Builder{}
		})

		JustBeforeEach(func() {
			err = cfg.Write(writer)
			actual = writer.String()
		})

		Context("with empty config", func() {
			BeforeEach(func() {
				cfg = &config.Config{}
				expected = ""
			})

			It("should not error", func() {
				Expect(err).To(MatchError("unable to write empty config"))
			})

			It("should return an empty config", func() {
				Expect(actual).To(Equal(expected))
			})
		})

		Context("with simple config", func() {
			BeforeEach(func() {
				cfg = &config.Config{
					Templates: []config.Template{
						{
							Repository: "https://github.com/rogueserenity/stenciler-test",
							Directory:  "test",
						},
					},
				}
				expected = `templates:
    - repository: https://github.com/rogueserenity/stenciler-test
      directory: test
`
			})

			It("should not error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should return an empty config", func() {
				Expect(actual).To(Equal(expected))
			})
		})

		Context("with full config", func() {
			BeforeEach(func() {
				cfg = &config.Config{
					Templates: []config.Template{
						{
							Repository: "https://github.com/rogueserenity/stenciler-test",
							Directory:  "test",
							Params: []*config.Param{
								{
									Name:           "spoo",
									Prompt:         "Would you like spoo?",
									Default:        "false",
									ValidationHook: "hooks/spoo.sh",
									Value:          "true",
								},
							},
							InitOnlyPaths: []string{
								"go.mod",
								"go.sum",
							},
							RawCopyPaths: []string{
								"**/*.py",
							},
							PreInitHookPaths: []string{
								"hooks/preinit",
							},
							PostInitHookPaths: []string{
								"hooks/postinit",
							},
							PreUpdateHookPaths: []string{
								"hooks/preup",
							},
							PostUpdateHookPaths: []string{
								"hooks/postup",
							},
						},
					},
				}
				expected = `templates:
    - repository: https://github.com/rogueserenity/stenciler-test
      directory: test
      params:
        - name: spoo
          prompt: Would you like spoo?
          default: "false"
          validation-hook: hooks/spoo.sh
          value: "true"
      init-only:
        - go.mod
        - go.sum
      raw-copy:
        - '**/*.py'
      pre-init-hooks:
        - hooks/preinit
      post-init-hooks:
        - hooks/postinit
      pre-update-hooks:
        - hooks/preup
      post-update-hooks:
        - hooks/postup
`
			})

			It("should not error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should return an empty config", func() {
				Expect(actual).To(Equal(expected))
			})
		})
	})

})
