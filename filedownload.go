package filedownload

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type urlInfo struct {
	url string
	fileLength int64
	fileName string
	etag string
	connSuccess bool
	statusCode int
}

type Chunk struct {
	Begin int64
	End int64
}

type Downloader struct {
	Timeout time.Duration
	fileName string
	fileSize int64
	chunk	[]Chunk
	chunkSize int64
	numberOfConnections int
	numberOfChunks int
	uri	string
}

type funcFeedback struct {
	id	int
	chunk	Chunk
	cursor int64
}

func NewDownloader() *Downloader {
	return &Downloader{
		Timeout: 15 * time.Second,
	}
}

func Download(uri string, nc int) error {
	d := &Downloader{
		uri:            uri,
		Timeout:        15 * time.Second,
		numberOfChunks: nc,
	}
	available := make(chan bool)
	done := make(chan bool)
	feedBack := make(chan *funcFeedback)

	downloadChunk := func(file *os.File, chunkIndex int) {
		for {
			<-available

		client := http.Client{}

		req, err := http.NewRequest("GET", d.uri, nil)
		if err != nil {
			return
		}

		req.Header.Add(
			"Range",
			fmt.Sprintf("bytes=%d-%d", d.chunk[chunkIndex].Begin, d.chunk[chunkIndex].End),
		)

		response, err := client.Do(req)
		defer func() {
			_ = response.Body.Close()
		}()

		buf := make([]byte, d.chunkSize+2)
		cursor := d.chunk[chunkIndex].Begin

		n, err := io.ReadFull(response.Body, buf)
		if err == io.EOF {
			log.Fatal(err)
		}

		_, errWr := file.WriteAt(buf[:n], cursor)
		if errWr != nil {
			log.Fatal(err)
			break
		}
		feedBack <- &funcFeedback{
			id:     chunkIndex,
			chunk:  d.chunk[chunkIndex],
			cursor: cursor,
		}

		done <- true
		break
	}
}

	info, err := gatherFileInfo(d)
	if err != nil {
		return err
	}
	d.fileName = info.fileName
	d.fileSize = info.fileLength

	_, errF := setupFile(d)
	if errF != nil {
		return err
	}

	file, err := os.OpenFile(d.fileName, os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	if !(info.connSuccess && info.statusCode == http.StatusOK) {
		return errors.New("cant get file info")
	}

	if err := createChunk(d); err != nil {
		return err
	}


	for i:=0 ; i < d.numberOfChunks; i++ {
		go downloadChunk(file, i)
		available <- true
	}

	count := nc

	r: for {
		select {
			case <-done:
				count--
				if count == 0 {
					break r
				}
			case f:= <- feedBack:
				fmt.Printf("%+v\n", f)
		}
	}
	return nil
}