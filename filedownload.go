package filedownload

import (
	"errors"
	"net/http"
	"os"
	"path"
	"time"
)

type urlInfo struct {
	url         string
	fileLength  int64
	fileName    string
	etag        string
	connSuccess bool
	statusCode  int
}

type Chunk struct {
	Begin int64
	End   int64
}

type Downloader struct {
	Timeout             time.Duration
	dir                 string
	fileName            string
	fileSize            int64
	chunk               []Chunk
	chunkSize           int64
	numberOfConnections int
	numberOfChunks      int64
	uri                 string
	info                urlInfo
}
type status struct {
	success bool
	error   error
}

type funcFeedback struct {
	id     int64
	chunk  Chunk
	cursor int64
	stat   status
}

func Download(uri string, nConn int, timeOut time.Duration, dir string) error {
	var feedbacks []*funcFeedback

	d := &Downloader{
		uri:                 uri,
		Timeout:             timeOut,
		numberOfConnections: nConn,
		dir:                 dir,
	}
	available := make(chan int64, 10)
	//done := make(chan bool)
	feedBack := make(chan *funcFeedback)

	info, err := gatherFileInfo(d)
	if err != nil {
		return err
	}
	//d.fileName = info.fileName
	//d.fileSize = info.fileLength

	_, errF := setupFile(d)
	if errF != nil {
		return err
	}

	file, err := os.OpenFile(path.Join(d.dir, d.fileName), os.O_WRONLY, 0666)
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

	for i := 0; i < d.numberOfConnections; i++ {
		go downloadChunk(d, file, available, feedBack)
	}

	//for chIdx:= int64(0); chIdx < d.numberOfChunks; chIdx++ {
	//	available <- chIdx
	//}

	count := d.numberOfChunks
	var chIdx int64 = 0
r:
	for {
		if chIdx < d.numberOfChunks {
			available <- chIdx
			chIdx++
		}
		select {
		case f := <-feedBack:
			count--
			if count == 0 {
				break r
			}
			feedbacks = append(feedbacks, f)
			//fmt.Printf("chunk %d stat: %+v error: %v\n", f.id, f.stat.success, f.stat.error)
		}
	}

	//check if chunk download fail run again
	for _, feedback := range feedbacks {
		if !feedback.stat.success {
			//go downloadChunk(d, file, feedback.id, available, feedBack)
			available <- feedback.id
		}
	}
	return nil
}
