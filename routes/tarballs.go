package routes

import (
	"fmt"
	"github.com/gillesdemey/npm-registry/storage"
	"io"
	"log"
	"net/http"
)

// GetTarball fetches a tarball from the upstream registry and falls back
// to storage engine if it fails
func GetTarball(w http.ResponseWriter, req *http.Request) {
	filename := req.URL.Query().Get(":filename")
	pkg := req.URL.Query().Get(":pkg")
	storage := StorageFromContext(req.Context())

	log.Printf("Fetching tarball %s/%s...", pkg, filename)
	resp, err := tryUpstreamTarball(pkg, filename)
	if err != nil {
		err := tryStorageTarball(storage, pkg, filename, w)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	if resp.StatusCode == http.StatusNotModified {
		return
	}

	// tee duplicates the body and writes to the ResponseWriter
	tee := io.TeeReader(resp.Body, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	updateTarballStorage(storage, pkg, filename, tee)
}

func tryUpstreamTarball(pkg string, filename string) (*http.Response, error) {
	log.Printf("Trying tarball upstream for %s/%s", pkg, filename)

	pkgMetaURL := fmt.Sprintf("https://registry.npmjs.org/%s/-/%s", pkg, filename)
	response, err := http.Get(pkgMetaURL)
	if err != nil {
		log.Printf("Failed to try upstream: %s", err)
		return nil, err
	}

	if response.StatusCode == http.StatusNotFound {
		err := fmt.Errorf("no such package available")
		log.Printf("Failed to try upstream: %s", err)
		return nil, err
	}

	// TODO catch 5xx errors from upstream

	return response, nil
}

func tryStorageTarball(s storage.Engine, pkg string, filename string, writer io.Writer) error {
	log.Printf("Trying storage tarball for %s/%s", pkg, filename)
	err := s.RetrieveTarball(pkg, filename, writer)
	if err != nil {
		log.Printf("Failed to try tarball storage: %s", err)
		return err
	}
	return nil
}

func updateTarballStorage(s storage.Engine, pkg string, filename string, data io.Reader) error {
	log.Printf("Updating %s/%s in tarball storage", pkg, filename)

	err := s.StoreTarball(pkg, filename, data)
	if err != nil {
		log.Printf("Failed to update %s/%s in tarball storage: %s", pkg, filename, err)
		return err
	}
	return nil
}
