package files_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/rogueserenity/stenciler/config"
	"github.com/rogueserenity/stenciler/files"
)

// since we're doing tests on the filesystem, keep the executions to a minimum for performance

var _ = Describe("CopyRaw", func() {

	var (
		template *config.Template
		srcDir   string
		destDir  string
		err      error
	)

	BeforeEach(func() {
		srcDir, err = os.MkdirTemp("", "copy-raw-test-src")
		Expect(err).ToNot(HaveOccurred())
		destDir, err = os.MkdirTemp("", "copy-raw-test-dst")
		Expect(err).ToNot(HaveOccurred())

		err = os.Chdir(destDir)
		Expect(err).ToNot(HaveOccurred())
		err = os.MkdirAll(srcDir+"/root/foo/bar/baz", 0755)
		Expect(err).ToNot(HaveOccurred())
		err = os.WriteFile(srcDir+"/root/foo.md", []byte("foo"), 0644)
		Expect(err).ToNot(HaveOccurred())
		err = os.WriteFile(srcDir+"/root/foo/foo.md", []byte("foo"), 0644)
		Expect(err).ToNot(HaveOccurred())
		err = os.WriteFile(srcDir+"/root/bar.txt", []byte("bar"), 0644)
		Expect(err).ToNot(HaveOccurred())
		err = os.WriteFile(srcDir+"/root/foo/bar.txt", []byte("bar"), 0644)
		Expect(err).ToNot(HaveOccurred())
		err = os.WriteFile(srcDir+"/root/foo/bar/bar.txt", []byte("bar"), 0644)
		Expect(err).ToNot(HaveOccurred())
		err = os.WriteFile(srcDir+"/root/foo/bar/baz/bar.txt", []byte("bar"), 0644)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		os.RemoveAll(srcDir)
		os.RemoveAll(destDir)
	})

	JustBeforeEach(func() {
		err = files.CopyRaw(srcDir, template)
	})

	Context("copies nothing with empty RawCopyPaths", func() {
		BeforeEach(func() {
			template = &config.Template{
				Directory: "root",
			}
		})

		It("should not error or copy any files", func() {
			Expect(err).ToNot(HaveOccurred())
			_, lerr := os.Stat(destDir + "/foo.md")
			Expect(lerr).To(HaveOccurred())
			_, lerr = os.Stat(destDir + "foo/foo.md")
			Expect(lerr).To(HaveOccurred())
			_, lerr = os.Stat(destDir + "/bar.txt")
			Expect(lerr).To(HaveOccurred())
			_, lerr = os.Stat(destDir + "/foo/bar.txt")
			Expect(lerr).To(HaveOccurred())
			_, lerr = os.Stat(destDir + "/foo/bar/bar.txt")
			Expect(lerr).To(HaveOccurred())
			_, lerr = os.Stat(destDir + "/foo/bar/baz/bar.txt")
			Expect(lerr).To(HaveOccurred())
		})
	})

	Context("copies explicit file names", func() {
		BeforeEach(func() {
			template = &config.Template{
				Directory: "root",
				RawCopyPaths: []string{
					"foo.md",
					"bar.txt",
					"foo/bar.txt",
				},
			}
		})

		It("should not error and should the copy files", func() {
			Expect(err).ToNot(HaveOccurred())
			_, lerr := os.Stat(destDir + "/foo.md")
			Expect(lerr).ToNot(HaveOccurred())
			_, lerr = os.Stat(destDir + "foo/foo.md")
			Expect(lerr).To(HaveOccurred())
			_, lerr = os.Stat(destDir + "/bar.txt")
			Expect(lerr).ToNot(HaveOccurred())
			_, lerr = os.Stat(destDir + "/foo/bar.txt")
			Expect(lerr).ToNot(HaveOccurred())
			_, lerr = os.Stat(destDir + "/foo/bar/bar.txt")
			Expect(lerr).To(HaveOccurred())
			_, lerr = os.Stat(destDir + "/foo/bar/baz/bar.txt")
			Expect(lerr).To(HaveOccurred())
		})
	})

	Context("copies globbed file names", func() {
		BeforeEach(func() {
			template = &config.Template{
				Directory: "root",
				RawCopyPaths: []string{
					"foo.md",
					"**/bar.txt",
				},
			}
		})

		It("should not error and should the copy files", func() {
			Expect(err).ToNot(HaveOccurred())
			_, lerr := os.Stat(destDir + "/foo.md")
			Expect(lerr).ToNot(HaveOccurred())
			_, lerr = os.Stat(destDir + "foo/foo.md")
			Expect(lerr).To(HaveOccurred())
			_, lerr = os.Stat(destDir + "/bar.txt")
			Expect(lerr).ToNot(HaveOccurred())
			_, lerr = os.Stat(destDir + "/foo/bar.txt")
			Expect(lerr).ToNot(HaveOccurred())
			_, lerr = os.Stat(destDir + "/foo/bar/bar.txt")
			Expect(lerr).ToNot(HaveOccurred())
			_, lerr = os.Stat(destDir + "/foo/bar/baz/bar.txt")
			Expect(lerr).ToNot(HaveOccurred())
		})
	})
})
