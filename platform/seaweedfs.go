package platform

import (
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/codern-org/codern/internal/constant"
	"github.com/linxGnu/goseaweedfs"
)

type SeaweedFs struct {
	client *goseaweedfs.Seaweed
}

func NewSeaweedFs(url string, filerUrl string) (*SeaweedFs, error) {
	httpClient := &http.Client{Timeout: 1 * time.Minute}
	client, err := goseaweedfs.NewSeaweed(
		url, []string{filerUrl},
		int64(constant.SeaweedFsChunkSize),
		httpClient,
	)
	if err != nil {
		return nil, err
	}
	return &SeaweedFs{
		client: client,
	}, err
}

func (fs *SeaweedFs) Upload(content io.Reader, size int, path string) error {
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

func (fs *SeaweedFs) Close() {
	fs.client.Close()
}
