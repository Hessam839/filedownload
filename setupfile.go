package filedownload

import (
	"os"
	"path/filepath"
)

func setupFile(d *Downloader) (os.FileInfo, error) {
	file, err := os.Create(
		filepath.Join(d.dir,
			d.fileName))

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
