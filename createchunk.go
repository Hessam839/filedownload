package filedownload

import "fmt"

func createChunk(d *Downloader) error {
	n := int64(d.numberOfChunks)
	reminder := d.fileSize % n
	roundedFileSize := d.fileSize - reminder
	chunkSize := roundedFileSize /n

	d.chunkSize = chunkSize

	d.chunk = make([]Chunk, d.numberOfChunks)

	boundary := int64(0)
	nextBoundary := chunkSize

	for i := int64(0); i < n; i++ {
		if reminder > 0 {
			nextBoundary++
			reminder--
		}
		d.chunk[i] = Chunk{Begin: boundary, End: nextBoundary}
		boundary = nextBoundary
		nextBoundary = boundary + chunkSize
	}

	fmt.Print(chunkSize)
	return nil
}
