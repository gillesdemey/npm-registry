package routes

import (
	"bytes"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gillesdemey/npm-registry/storage"
	"io"
	"net/http"
)

// GetTarball fetches a tarball from the upstream registry and falls back
// to storage engine if it fails
func GetTarball(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	storage := StorageFromContext(req.Context())

	filename := req.URL.Query().Get(":filename")
	scope := req.URL.Query().Get(":scope")
	pkg := req.URL.Query().Get(":pkg")

	if scope != "" && pkg != "" {
		pkg = scope + "/" + pkg
	}

	buff := new(bytes.Buffer)
	multiWriter := io.MultiWriter(w, buff)

	resp, err := tryUpstreamTarball(pkg, filename)
	if err != nil {
		err := tryStorageTarball(storage, pkg, filename, multiWriter)
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

	_, err = io.Copy(multiWriter, resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if scope != "" && filename != "" {
		filename = scope + "/" + filename
	}
	updateTarballStorage(storage, filename, buff)
}

func tryUpstreamTarball(pkg string, filename string) (*http.Response, error) {
	logger := log.WithFields(log.Fields{
		"package": pkg,
		"tarball": filename,
		"source":  "upstream",
	})
	logger.Info("Trying upstream...")

	pkgMetaURL := fmt.Sprintf("https://registry.npmjs.org/%s/-/%s", pkg, filename)
	response, err := http.Get(pkgMetaURL)
	if err != nil {
		logger.Warn("Upstream failed: ", err)
		return nil, err
	}

	if response.StatusCode == http.StatusNotFound {
		err := errors.New("no such package available")
		logger.Warn("Upstream failed: ", err)
		return nil, err
	}

	// TODO catch 5xx errors from upstream

	return response, nil
}

func tryStorageTarball(s storage.Engine, pkg string, filename string, writer io.Writer) error {
	logger := log.WithFields(log.Fields{
		"package": pkg,
		"tarball": filename,
		"source":  "upstream",
	})
	logger.Info("Trying storage...")

	err := s.RetrieveTarball(pkg, filename, writer)
	if err != nil {
		logger.Warn("no such package available")
		return err
	}
	return nil
}

func updateTarballStorage(s storage.Engine, filename string, data io.Reader) error {
	logger := log.WithFields(log.Fields{
		"tarball": filename,
		"source":  "upstream",
	})
	logger.Info("Updating tarball storage")

	err := s.StoreTarball(filename, data)
	if err != nil {
		logger.Error("Failed to update tarball storage: %s", err)
		return err
	}
	return nil
}
