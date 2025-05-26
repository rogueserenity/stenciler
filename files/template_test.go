package files_test

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/rogueserenity/stenciler/config"
	"github.com/rogueserenity/stenciler/files"
)

type CopyTemplatedTestSuite struct {
	suite.Suite

	srcDir  string
	destDir string
}

func TestCopyTemplatedTestSuite(t *testing.T) {
	suite.Run(t, new(CopyTemplatedTestSuite))
}

func (s *CopyTemplatedTestSuite) SetupSuite() {
	var err error
	s.srcDir, err = os.MkdirTemp("", "templated-test-src")
	s.Require().NoError(err)

	err = os.MkdirAll(path.Join(s.srcDir, "/root/foo/bar/baz"), 0755)
	s.Require().NoError(err)
	err = os.WriteFile(path.Join(s.srcDir, "/root/foo.md"), []byte("{{.foo}}"), 0644)
	s.Require().NoError(err)
	err = os.WriteFile(path.Join(s.srcDir, "/root/foo/foo.md"), []byte("{{.foo}}"), 0644)
	s.Require().NoError(err)
	err = os.WriteFile(path.Join(s.srcDir, "/root/bar.txt"), []byte("{{.bar}}"), 0644)
	s.Require().NoError(err)
	err = os.WriteFile(path.Join(s.srcDir, "/root/foo/bar.txt"), []byte("{{.bar}}"), 0644)
	s.Require().NoError(err)
	err = os.WriteFile(path.Join(s.srcDir, "/root/foo/bar/bar.txt"), []byte("{{.bar}}"), 0644)
	s.Require().NoError(err)
	err = os.WriteFile(path.Join(s.srcDir, "/root/foo/bar/baz/bar.txt"), []byte("{{.bar}}"), 0644)
	s.Require().NoError(err)
}

func (s *CopyTemplatedTestSuite) TearDownSuite() {
	os.RemoveAll(s.srcDir)
}

func (s *CopyTemplatedTestSuite) SetupTest() {
	var err error
	s.destDir, err = os.MkdirTemp("", "templated-test-dst")
	s.Require().NoError(err)

	err = os.Chdir(s.destDir)
	s.Require().NoError(err)
}

func (s *CopyTemplatedTestSuite) TearDownTest() {
	os.RemoveAll(s.destDir)
}

func (s *CopyTemplatedTestSuite) TestCopyTemplated() {
	template := &config.Template{
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

	err := files.CopyTemplated(s.srcDir, template)
	s.Require().NoError(err)

	b, err := os.ReadFile(path.Join(s.destDir, "/foo.md"))
	s.Require().NoError(err)
	s.Equal("foo", string(b))

	b, err = os.ReadFile(path.Join(s.destDir, "/foo/foo.md"))
	s.Require().NoError(err)
	s.Equal("foo", string(b))

	b, err = os.ReadFile(path.Join(s.destDir, "/bar.txt"))
	s.Require().NoError(err)
	s.Equal("bar", string(b))

	b, err = os.ReadFile(path.Join(s.destDir, "/foo/bar.txt"))
	s.Require().NoError(err)
	s.Equal("bar", string(b))

	b, err = os.ReadFile(path.Join(s.destDir, "/foo/bar/bar.txt"))
	s.Require().NoError(err)
	s.Equal("bar", string(b))

	b, err = os.ReadFile(path.Join(s.destDir, "/foo/bar/baz/bar.txt"))
	s.Require().NoError(err)
	s.Equal("bar", string(b))
}
