package routes

import (
	"fmt"
	"io"
	"log"
	"net/http"

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

	resp, err := tryUpstream(pkg)
	if err != nil {
		err = tryMetaStorage(storage, pkg, w)
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

	tee := io.TeeReader(resp.Body, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	updateMetaStorage(storage, pkg, tee)
}

func tryUpstream(pkg string) (*http.Response, error) {
	log.Printf("Trying upstream for %s\n", pkg)

	pkgMetaURL := fmt.Sprintf("https://registry.npmjs.org/%s", pkg)
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

func tryMetaStorage(s storage.Engine, pkg string, writer io.Writer) error {
	log.Printf("Trying storage for %s\n", pkg)

	err := s.RetrieveMetadata(pkg, writer)
	if err != nil {
		log.Printf("Failed to try storage: %s", err)
		return err
	}

	return nil
}

func updateMetaStorage(s storage.Engine, pkg string, data io.Reader) error {
	log.Printf("Updating %s in meta storage", pkg)

	err := s.StoreMetadata(pkg, data)
	if err != nil {
		log.Printf("Failed to update %s in meta storage: %s", pkg, err)
		return err
	}

	return nil
}
