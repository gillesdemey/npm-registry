package packages

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestRewriteTarballLocation(t *testing.T) {
	reader, _ := os.Open("fixtures/express.json")
	err := RewriteTarballLocation(reader, ioutil.Discard)

	if err != nil {
		t.Error(err)
	}
}

func BenchmarkRewriteTarballLocation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		reader, _ := os.Open("fixtures/express.json")
		RewriteTarballLocation(reader, ioutil.Discard)
	}
}
