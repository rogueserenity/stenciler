package files_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/rogueserenity/stenciler/config"
	"github.com/rogueserenity/stenciler/files"
)

// since we're doing tests on the filesystem, keep the executions to a minimum for performance

var _ = Describe("CopyTemplated", func() {

	var (
		template *config.Template
		srcDir   string
		destDir  string
		err      error
	)

	BeforeEach(func() {
		srcDir, err = os.MkdirTemp("", "templated-test-src")
		Expect(err).ToNot(HaveOccurred())
		destDir, err = os.MkdirTemp("", "templated-test-dst")
		Expect(err).ToNot(HaveOccurred())

		err = os.Chdir(destDir)
		Expect(err).ToNot(HaveOccurred())
		err = os.MkdirAll(srcDir+"/root/foo/bar/baz", 0755)
		Expect(err).ToNot(HaveOccurred())
		err = os.WriteFile(srcDir+"/root/foo.md", []byte("{{.foo}}"), 0644)
		Expect(err).ToNot(HaveOccurred())
		err = os.WriteFile(srcDir+"/root/foo/foo.md", []byte("{{.foo}}"), 0644)
		Expect(err).ToNot(HaveOccurred())
		err = os.WriteFile(srcDir+"/root/bar.txt", []byte("{{.bar}}"), 0644)
		Expect(err).ToNot(HaveOccurred())
		err = os.WriteFile(srcDir+"/root/foo/bar.txt", []byte("{{.bar}}"), 0644)
		Expect(err).ToNot(HaveOccurred())
		err = os.WriteFile(srcDir+"/root/foo/bar/bar.txt", []byte("{{.bar}}"), 0644)
		Expect(err).ToNot(HaveOccurred())
		err = os.WriteFile(srcDir+"/root/foo/bar/baz/bar.txt", []byte("{{.bar}}"), 0644)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		os.RemoveAll(srcDir)
		os.RemoveAll(destDir)
	})

	JustBeforeEach(func() {
		err = files.CopyTemplated(srcDir, template)
	})

	Context("copies globbed file names", func() {
		BeforeEach(func() {
			template = &config.Template{
				Directory: "root",
				Params: []*config.Param{
					{
						Name:  "foo",
						Value: "foo",
					},
					{
						Name:  "bar",
						Value: "bar",
					},
				},
			}
		})

		It("should not error and should the copy files", func() {
			Expect(err).ToNot(HaveOccurred())
			b, lerr := os.ReadFile(destDir + "/foo.md")
			Expect(lerr).ToNot(HaveOccurred())
			Expect(string(b)).To(Equal("foo"))
			b, lerr = os.ReadFile(destDir + "/foo/foo.md")
			Expect(lerr).ToNot(HaveOccurred())
			Expect(string(b)).To(Equal("foo"))
			b, lerr = os.ReadFile(destDir + "/bar.txt")
			Expect(lerr).ToNot(HaveOccurred())
			Expect(string(b)).To(Equal("bar"))
			b, lerr = os.ReadFile(destDir + "/foo/bar.txt")
			Expect(lerr).ToNot(HaveOccurred())
			Expect(string(b)).To(Equal("bar"))
			b, lerr = os.ReadFile(destDir + "/foo/bar/bar.txt")
			Expect(lerr).ToNot(HaveOccurred())
			Expect(string(b)).To(Equal("bar"))
			b, lerr = os.ReadFile(destDir + "/foo/bar/baz/bar.txt")
			Expect(lerr).ToNot(HaveOccurred())
			Expect(string(b)).To(Equal("bar"))
		})
	})
})
