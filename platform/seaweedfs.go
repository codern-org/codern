package platform

import (
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/codern-org/codern/domain"
	"github.com/linxGnu/goseaweedfs"
)

type seaweedFs struct {
	client *goseaweedfs.Seaweed
}

func NewSeaweedFs(url string, filerUrl []string) (domain.SeaweedFs, error) {
	httpClient := &http.Client{Timeout: 1 * time.Minute}
	client, err := goseaweedfs.NewSeaweed(url, filerUrl, 4096, httpClient)
	if err != nil {
		return nil, err
	}
	return &seaweedFs{
		client: client,
	}, err
}

func (fs *seaweedFs) Upload(content io.Reader, size int, path string) error {
	filer := fs.client.Filers()[0]
	if filer == nil {
		return errors.New("cannot connect to file system upstream")
	}
	_, err := filer.Upload(content, int64(size), path, "", "")
	if err != nil {
		return err
	}
	return nil
}
