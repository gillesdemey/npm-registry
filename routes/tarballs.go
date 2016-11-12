package routes

import (
	"io"
	"log"
	"net/http"
)

func GetTarball(w http.ResponseWriter, req *http.Request) {
	filename := req.URL.Query().Get(":filename")
	storage := StorageFromContext(req.Context())

	log.Printf("Fetching tarball %s...", filename)
	err := storage.RetrieveTarball(filename, w)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func tryStorageTarball(pkg string) (io.Reader, error) {
	return nil, nil
}

func tryUpstreamTarball(pkg string) (*http.Response, error) {
	return nil, nil
}
