package routes

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/Jeffail/gabs"
	log "github.com/Sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

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

	logger.Info("Trying storage...")

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

// TODO add authentication middleware
func PublishPackage(w http.ResponseWriter, req *http.Request) {
	var err error
	storage := StorageFromContext(req.Context())
	renderer := RendererFromContext(req.Context())

	pkgInfo, err := gabs.ParseJSONBuffer(req.Body)
	if err != nil {
		renderer.JSON(w, http.StatusBadRequest, map[string]string{"error": "bad request"})
		return
	}
	defer req.Body.Close()

	pkgName, _ := pkgInfo.Path("name").Data().(string)

	// Check if we have 1 dist-tag
	distTags, _ := pkgInfo.Path("dist-tags").ChildrenMap()
	if len(distTags) != 1 {
		renderer.JSON(w, http.StatusBadRequest, map[string]string{"error": "must have 1 dist-tag"})
		return
	}

	// Check if version of this package already exists
	var tag string
	for k, _ := range distTags {
		tag = k
		break
	}
	version := distTags[tag].Data().(string)

	logger := log.WithFields(log.Fields{
		"package": pkgName,
		"version": version,
	})
	logger.Info("Processing package")

	err = storage.RetrieveTarball(
		pkgName,
		// TODO there's probably a better way then manually creating the filename?
		fmt.Sprintf("%s-%s.tgz", pkgName, version),
		ioutil.Discard, // write to /dev/null
	)

	if err == nil {
		renderer.JSON(w, http.StatusBadRequest, map[string]string{
			"error": "version already exists",
		})
		return
	}

	attachments, _ := pkgInfo.Path("_attachments").ChildrenMap()
	for filename, attachment := range attachments {
		logger.WithFields(log.Fields{"attachment": filename}).Info("Decoding...")
		data := attachment.Path("data").Data().(string)
		buff, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		storage.StoreTarball(pkgName, filename, bytes.NewReader(buff))
	}

	// TODO save package metadata
	// 1. check if package metadata already exists
	// 2. add dist-tag version to existing versions
	// 3. delete _attachments from JSON payload
	pkgInfo.Delete("_attachments")
	storage.StoreMetadata(pkgName, strings.NewReader(pkgInfo.String()))
}
