package prompt_test

import (
	"bytes"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/rogueserenity/stenciler/config"
	"github.com/rogueserenity/stenciler/prompt"
)

var _ = Describe("ForParamValuesWithInOut", func() {
	var (
		template *config.Template
		repoDir  string
		stdin    bytes.Buffer
		stdout   *strings.Builder

		err error
	)

	BeforeEach(func() {
		stdout = &strings.Builder{}
	})

	JustBeforeEach(func() {
		err = prompt.ForParamValuesWithInOut(template, repoDir, &stdin, stdout)
	})

	Context("with a param with no prompt", func() {
		BeforeEach(func() {
			template = &config.Template{
				Params: []*config.Param{
					{
						Name:  "foo",
						Value: "bar",
					},
				},
			}
		})

		It("should not error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("should not prompt for a value", func() {
			Expect(stdout.String()).To(BeEmpty())
		})
	})

	Context("with a param with a prompt", func() {
		BeforeEach(func() {
			template = &config.Template{
				Params: []*config.Param{
					{
						Name:   "foo",
						Prompt: "Enter a value for foo",
					},
				},
			}
		})

		Context("with a value", func() {
			BeforeEach(func() {
				template.Params[0].Value = "bar"
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should not prompt for a value", func() {
				Expect(stdout.String()).To(BeEmpty())
			})
		})

		Context("without a value", func() {
			Context("with a default value", func() {
				BeforeEach(func() {
					template.Params[0].Default = "baz"
				})

				Context("without input", func() {
					BeforeEach(func() {
						stdin.WriteString("\n")
					})

					It("should not error", func() {
						Expect(err).NotTo(HaveOccurred())
					})

					It("should prompt for a value", func() {
						Expect(stdout.String()).To(Equal("Enter a value for foo [baz]: "))
					})

					It("should set the value to the default value", func() {
						Expect(template.Params[0].Value).To(Equal("baz"))
					})
				})

				Context("with input", func() {
					BeforeEach(func() {
						stdin.WriteString("blah\n")
					})

					It("should not error", func() {
						Expect(err).NotTo(HaveOccurred())
					})

					It("should prompt for a value", func() {
						Expect(stdout.String()).To(Equal("Enter a value for foo [baz]: "))
					})

					It("should set the value to the input value", func() {
						Expect(template.Params[0].Value).To(Equal("blah"))
					})
				})
			})

			Context("without a default value", func() {
				BeforeEach(func() {
					stdin.WriteString("blah\n")
				})

				It("should not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should prompt for a value", func() {
					Expect(stdout.String()).To(Equal("Enter a value for foo: "))
				})

				It("should set the value to the input value", func() {
					Expect(template.Params[0].Value).To(Equal("blah"))
				})
			})
		})
	})
})
