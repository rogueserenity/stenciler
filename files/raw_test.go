package files_test

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/rogueserenity/stenciler/config"
	"github.com/rogueserenity/stenciler/files"
)

type CopyRawTestSuite struct {
	suite.Suite

	srcDir  string
	destDir string
}

func TestCopyRawTestSuite(t *testing.T) {
	suite.Run(t, new(CopyRawTestSuite))
}

func (s *CopyRawTestSuite) SetupSuite() {
	var err error
	s.srcDir, err = os.MkdirTemp("", "copy-raw-test-src")
	s.Require().NoError(err)

	err = os.MkdirAll(path.Join(s.srcDir, "/root/foo/bar/baz"), 0755)
	s.Require().NoError(err)
	err = os.WriteFile(path.Join(s.srcDir, "/root/foo.md"), []byte("foo"), 0644)
	s.Require().NoError(err)
	err = os.WriteFile(path.Join(s.srcDir, "/root/foo/foo.md"), []byte("foo"), 0644)
	s.Require().NoError(err)
	err = os.WriteFile(path.Join(s.srcDir, "/root/bar.txt"), []byte("bar"), 0644)
	s.Require().NoError(err)
	err = os.WriteFile(path.Join(s.srcDir, "/root/foo/bar.txt"), []byte("bar"), 0644)
	s.Require().NoError(err)
	err = os.WriteFile(path.Join(s.srcDir, "/root/foo/bar/bar.txt"), []byte("bar"), 0644)
	s.Require().NoError(err)
	err = os.WriteFile(path.Join(s.srcDir, "/root/foo/bar/baz/bar.txt"), []byte("bar"), 0644)
	s.Require().NoError(err)
}

func (s *CopyRawTestSuite) TearDownSuite() {
	os.RemoveAll(s.srcDir)
}

func (s *CopyRawTestSuite) SetupTest() {
	var err error
	s.destDir, err = os.MkdirTemp("", "copy-raw-test-dst")
	s.Require().NoError(err)

	err = os.Chdir(s.destDir)
	s.Require().NoError(err)
}

func (s *CopyRawTestSuite) TearDownTest() {
	os.RemoveAll(s.destDir)
}

func (s *CopyRawTestSuite) TestEmptyRawCopyPaths() {
	template := &config.Template{
		Directory: "root",
	}

	err := files.CopyRaw(s.srcDir, template)
	s.Require().NoError(err)

	s.Require().NoFileExists(path.Join(s.destDir, "foo.md"))
	s.Require().NoFileExists(path.Join(s.destDir, "foo/foo.md"))
	s.Require().NoFileExists(path.Join(s.destDir, "bar.txt"))
	s.Require().NoFileExists(path.Join(s.destDir, "foo/bar.txt"))
	s.Require().NoFileExists(path.Join(s.destDir, "foo/bar/bar.txt"))
	s.Require().NoFileExists(path.Join(s.destDir, "foo/bar/baz/bar.txt"))
}

func (s *CopyRawTestSuite) TestCopyExplicitFileNames() {
	template := &config.Template{
		Directory: "root",
		RawCopyPaths: []string{
			"foo.md",
			"bar.txt",
			"foo/bar.txt",
		},
	}

	err := files.CopyRaw(s.srcDir, template)
	s.Require().NoError(err)
	s.Require().FileExists(path.Join(s.destDir, "foo.md"))
	s.Require().NoFileExists(path.Join(s.destDir, "foo/foo.md"))
	s.Require().FileExists(path.Join(s.destDir, "bar.txt"))
	s.Require().FileExists(path.Join(s.destDir, "foo/bar.txt"))
	s.Require().NoFileExists(path.Join(s.destDir, "foo/bar/bar.txt"))
	s.Require().NoFileExists(path.Join(s.destDir, "foo/bar/baz/bar.txt"))
}

func (s *CopyRawTestSuite) TestCopyGlobbedFileNames() {
	template := &config.Template{
		Directory: "root",
		RawCopyPaths: []string{
			"foo.md",
			"**/bar.txt",
		},
	}

	err := files.CopyRaw(s.srcDir, template)
	s.Require().NoError(err)
	s.Require().FileExists(path.Join(s.destDir, "foo.md"))
	s.Require().NoFileExists(path.Join(s.destDir, "foo/foo.md"))
	s.Require().FileExists(path.Join(s.destDir, "bar.txt"))
	s.Require().FileExists(path.Join(s.destDir, "foo/bar.txt"))
	s.Require().FileExists(path.Join(s.destDir, "foo/bar/bar.txt"))
	s.Require().FileExists(path.Join(s.destDir, "foo/bar/baz/bar.txt"))
}
