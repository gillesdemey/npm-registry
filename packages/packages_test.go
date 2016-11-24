package packages

import (
	"bytes"
	"github.com/Jeffail/gabs"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestRewriteTarballLocation(t *testing.T) {
	reader := strings.NewReader(`{
		"tarball":"https://registry.npmjs.org/foo.tgz",
		"tarball2":"https://registry.npmjs.org/bar.tgz"
	}`)
	writer := new(bytes.Buffer)
	err := RewriteTarballLocation(reader, writer)

	if err != nil {
		t.Error(err)
	}

	match := string(writer.Bytes()) == `{
		"tarball":"https://localhost:8080/foo.tgz",
		"tarball2":"https://localhost:8080/bar.tgz"
	}`

	if !match {
		t.Error("string incorrectly transformed")
	}
}

func BenchmarkRewriteTarballLocation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		reader, _ := os.Open("fixtures/express.json")
		RewriteTarballLocation(reader, ioutil.Discard)
	}
}

func TestSplitPackageNameUnscoped(t *testing.T) {
	scope, pkgName := SplitPackageName("foo")
	assert.Equal(t, scope, "")
	assert.Equal(t, pkgName, "foo")
}

func TestSplitPackageNameScoped(t *testing.T) {
	scope, pkgName := SplitPackageName("@foo/bar")
	assert.Equal(t, scope, "@foo")
	assert.Equal(t, pkgName, "bar")
}

func TestRewriteScopedTarballs(t *testing.T) {
	input := make(map[string]*gabs.Container)
	input["1.0.0"], _ = gabs.ParseJSON([]byte(`{
		"dist": {
			"tarball": "http://foo.bar/@foo/bar/-/@foo/bar-1.0.0.tgz"
		}
	}`))

	expected := make(map[string]*gabs.Container)
	expected["1.0.0"], _ = gabs.ParseJSON([]byte(`{
		"dist": {
			"tarball": "http://foo.bar/@foo/bar/-/bar-1.0.0.tgz"
		}
	}`))

	output := RewriteScopedTarballs("@foo/bar", input)

	assert.Equal(t, output, expected)
}
