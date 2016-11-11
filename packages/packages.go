package packages

import (
	"io"
	"io/ioutil"
	"strings"
)

func RewriteTarballLocation(meta io.Reader, writer io.Writer) error {
	replacer := strings.NewReplacer(
		"https://registry.npmjs.org",
		"https://localhost:8080",
	)
	buff, err := ioutil.ReadAll(meta)
	if err != nil {
		return err
	}

	replacer.WriteString(writer, string(buff))
	return nil
}
