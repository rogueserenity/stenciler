package prompt_test

import (
	"bytes"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/rogueserenity/stenciler/config"
	"github.com/rogueserenity/stenciler/prompt"
)

var _ = Describe("SelectTemplate", func() {
	var (
		templateDir string
		cfg         *config.Config
		stdin       bytes.Buffer
		stdout      *strings.Builder

		template *config.Template
		err      error
	)

	BeforeEach(func() {
		stdout = &strings.Builder{}
	})

	JustBeforeEach(func() {
		template, err = prompt.SelectTemplateWithInOut(templateDir, cfg, &stdin, stdout)
	})

	Context("when a template directory is specified", func() {
		BeforeEach(func() {
			templateDir = "foo"
		})

		Context("when the template directory is found in the config", func() {
			BeforeEach(func() {
				cfg = &config.Config{
					Templates: []*config.Template{
						{
							Directory: "foo",
						},
						{
							Directory: "bar",
						},
					},
				}
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return the template", func() {
				Expect(template).NotTo(BeNil())
				Expect(template.Directory).To(Equal("foo"))
			})
		})

		Context("when the template directory is not found in the config", func() {
			BeforeEach(func() {
				cfg = &config.Config{
					Templates: []*config.Template{
						{
							Directory: "bar",
						},
					},
				}
			})

			It("should error", func() {
				Expect(err).To(HaveOccurred())
			})

			It("should return nil", func() {
				Expect(template).To(BeNil())
			})
		})
	})

	Context("when no template directory is specified", func() {
		BeforeEach(func() {
			templateDir = ""
		})

		Context("when there is only one template in the config", func() {
			BeforeEach(func() {
				cfg = &config.Config{
					Templates: []*config.Template{
						{
							Directory: "foo",
						},
					},
				}
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return the template", func() {
				Expect(template).NotTo(BeNil())
				Expect(template.Directory).To(Equal("foo"))
			})
		})

		Context("when there are multiple templates in the config", func() {
			BeforeEach(func() {
				cfg = &config.Config{
					Templates: []*config.Template{
						{
							Directory: "foo",
						},
						{
							Directory: "bar",
						},
					},
				}
			})

			Context("when the user selects a valid template directory", func() {
				BeforeEach(func() {
					stdin.Reset()
					stdin.WriteString("foo\n")
				})

				It("should have the correct output", func() {
					Expect(stdout.String()).To(Equal("Available templates:\n" +
						">  bar\n" +
						">  foo\n" +
						"please specify the template directory to use: "))
				})

				It("should not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should return the template", func() {
					Expect(template).NotTo(BeNil())
					Expect(template.Directory).To(Equal("foo"))
				})
			})

			Context("when the user selects an invalid template directory", func() {
				BeforeEach(func() {
					stdin.Reset()
					stdin.WriteString("baz\nfoo\n")
				})

				It("should have the correct output", func() {
					Expect(stdout.String()).To(Equal("Available templates:\n" +
						">  bar\n" +
						">  foo\n" +
						"please specify the template directory to use: " +
						"Available templates:\n" +
						">  bar\n" +
						">  foo\n" +
						"please specify the template directory to use: "))
				})

				It("should not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should return the template", func() {
					Expect(template).NotTo(BeNil())
					Expect(template.Directory).To(Equal("foo"))
				})
			})
		})
	})
})
