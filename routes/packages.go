package routes

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io"
	"net/http"

	"github.com/gillesdemey/npm-registry/packages"
	"github.com/gillesdemey/npm-registry/storage"
)

// GetPackageMetadata fetches package meta data
//
// 1. fetch package metadata from NPM upstream
// 2. if 304 -> return 304 and exit
// 3. if 404 or error -> check storage package
// 4. if still not found -> 404
// 5. if all is well, update storage with newer metadata
func GetPackageMetadata(w http.ResponseWriter, req *http.Request) {
	var err error

	pkg := req.URL.Query().Get(":pkg")
	storage := StorageFromContext(req.Context())
	renderer := RendererFromContext(req.Context())

	pr, pw := io.Pipe()
	multiWriter := io.MultiWriter(w, pw)

	resp, err := tryUpstream(pkg)
	if err != nil {
		err = tryMetaStorage(storage, pkg, multiWriter)
		if err != nil {
			renderer.JSON(w, http.StatusNotFound, map[string]string{
				"error": "no such package available",
			})
			return
		}
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	if resp.StatusCode == http.StatusNotModified {
		return
	}

	go func() {
		packages.RewriteTarballLocation(resp.Body, multiWriter)
		pw.Close()
	}()

	updateMetaStorage(storage, pkg, pr)
}

func tryUpstream(pkg string) (*http.Response, error) {
	logger := log.WithFields(log.Fields{
		"source":  "upstream",
		"package": pkg,
	})

	logger.Info("Trying upstream...")

	pkgMetaURL := fmt.Sprintf("https://registry.npmjs.org/%s", pkg)
	response, err := http.Get(pkgMetaURL)
	if err != nil {
		logger.Warn("Upstream failed: ", err)
		return nil, err
	}

	if response.StatusCode == http.StatusNotFound {
		err := fmt.Errorf("no such package available")
		logger.Warn(err)
		return nil, err
	}

	// TODO catch 5xx errors from upstream

	return response, nil
}

func tryMetaStorage(s storage.Engine, pkg string, writer io.Writer) error {
	logger := log.WithFields(log.Fields{
		"source":  "storage",
		"package": pkg,
	})

	err := s.RetrieveMetadata(pkg, writer)
	if err != nil {
		logger.Warn("no such package available")
		return err
	}

	return nil
}

func updateMetaStorage(s storage.Engine, pkg string, data io.Reader) error {
	logger := log.WithFields(log.Fields{
		"source":  "storage",
		"package": pkg,
	})
	logger.Info("Updating meta storage")

	err := s.StoreMetadata(pkg, data)
	if err != nil {
		logger.Error("Failed to update meta storage: ", err)
		return err
	}

	return nil
}
