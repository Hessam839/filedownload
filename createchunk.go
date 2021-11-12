package filedownload

func createChunk(d *Downloader) error {
	calculateChunk(d)
	n := d.numberOfChunks
	reminder := d.fileSize % n
	roundedFileSize := d.fileSize - reminder
	chunkSize := roundedFileSize / n

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

	//fmt.Print(chunkSize)
	return nil
}

func calculateChunk(d *Downloader) {
	d.numberOfChunks = 5
	if 100_000 < d.fileSize && d.fileSize <= 1_000_000 {
		d.numberOfChunks = 20
	} else if 1_000_000 < d.fileSize {
		d.numberOfChunks = d.fileSize / 200_000
	}

}
