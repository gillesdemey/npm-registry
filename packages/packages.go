package packages

import (
	"github.com/Jeffail/gabs"
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

func RewriteScopedTarballs(pkgName string, versions map[string]*gabs.Container) map[string]*gabs.Container {
	scope, pkgName := SplitPackageName(pkgName)
	if scope == "" { // not a scope package
		return versions
	}

	for _, version := range versions {
		oldPath := version.Path("dist.tarball").Data().(string)
		newPath := strings.Replace(oldPath, "-/"+scope, "-", 1)
		version.SetP(newPath, "dist.tarball")
	}
	return versions
}

func SplitPackageName(pkg string) (scope string, pkgName string) {
	isScoped := strings.Contains(pkg, "@")

	if isScoped {
		split := strings.Split(pkg, "/")
		scope := split[0]
		pkgName := split[1]

		return scope, pkgName
	} else {
		return "", pkg
	}
}
