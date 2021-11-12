package filedownload

import (
	"os"
	"path/filepath"
)

func setupFile(d *Downloader) (os.FileInfo, error) {
	path := filepath.Join(d.dir, d.fileName)
	file, err := os.Create(path)

	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	if err := file.Truncate(d.fileSize); err != nil {
		return nil, err
	}

	return file.Stat()
}
