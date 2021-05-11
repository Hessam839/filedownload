package filedownload

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func downloadChunk(d *Downloader, file *os.File, chunkIndex int, available chan bool, feedBack chan *funcFeedback) {
	//for {
	<-available

	client := http.Client{}
	if d.Timeout > 0 {
		client.Timeout = d.Timeout
	}

	cursor := d.chunk[chunkIndex].Begin

	req, errReq := http.NewRequest("GET", d.uri, nil)
	if errReq != nil {
		feedBack <- &funcFeedback{
			id:     chunkIndex,
			chunk:  d.chunk[chunkIndex],
			cursor: cursor,
			stat: status{
				success: false,
				error:   errReq,
			},
		}
	}

	req.Header.Add(
		"Range",
		fmt.Sprintf("bytes=%d-%d", d.chunk[chunkIndex].Begin, d.chunk[chunkIndex].End),
	)

	response, errCli := client.Do(req)
	if errCli != nil {
		feedBack <- &funcFeedback{
			id:     chunkIndex,
			chunk:  d.chunk[chunkIndex],
			cursor: cursor,
			stat: status{
				success: false,
				error:   errCli,
			},
		}
	}

	defer func() {
		_ = response.Body.Close()
	}()

	io.NopCloser(response.Body)

	buf := make([]byte, d.chunkSize+2)

	n, errRd := io.ReadFull(response.Body, buf)
	if errRd == io.EOF {
		feedBack <- &funcFeedback{
			id:     chunkIndex,
			chunk:  d.chunk[chunkIndex],
			cursor: cursor,
			stat: status{
				success: false,
				error:   errRd,
			},
		}
	}

	_, errWr := file.WriteAt(buf[:n], cursor)
	if errWr != nil {
		feedBack <- &funcFeedback{
			id:     chunkIndex,
			chunk:  d.chunk[chunkIndex],
			cursor: cursor,
			stat: status{
				success: false,
				error:   errWr,
			},
		}
		//break
	}

	feedBack <- &funcFeedback{
		id:     chunkIndex,
		chunk:  d.chunk[chunkIndex],
		cursor: cursor,
		stat: status{
			success: true,
			error:   nil,
		},
	}

	//break
	//}
}
