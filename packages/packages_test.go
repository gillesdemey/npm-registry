package packages

import (
	"bytes"
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
