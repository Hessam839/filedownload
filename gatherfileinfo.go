package filedownload

import (
	"errors"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
)

var (
	ErrorNotFound   = errors.New("url not found")
	ErrorBadRequest = errors.New("bad request url")
)

func gatherFileInfo(d *Downloader) (*urlInfo, error) {
	parsedUrl, err := url.Parse(d.uri)
	if err != nil {
		return nil, err
	}

	fileName := filepath.Base(parsedUrl.Path)
	client := http.Client{Timeout: d.Timeout}

	resp, err := client.Head(d.uri)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return nil, ErrorNotFound
		case http.StatusBadRequest:
			return nil, ErrorBadRequest
		}
	}
	fileSize, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 0, 64)
	if err != nil {
		return nil, err
	}
	etag := resp.Header.Get("Etag")

	d.fileName = fileName
	d.fileSize = fileSize

	return &urlInfo{
		statusCode:  resp.StatusCode,
		fileLength:  fileSize,
		fileName:    fileName,
		etag:        etag,
		url:         d.uri,
		connSuccess: true,
	}, nil
}
