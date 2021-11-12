package filedownload

import (
	"testing"
	"time"
)

func Test_GetFileInfo(t *testing.T) {
	d := NewDownloader()
	d.uri = `https://www.pezeshkonline.ir/download/doctormaleki-200519143319.jpg`
	info, err := gatherFileInfo(d)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("url info: %+v", info)
}

func Benchmark_GetFileInfo(b *testing.B) {
	b.ReportAllocs()
	d := NewDownloader()
	d.uri = `https://www.pezeshkonline.ir/download/doctormaleki-200519143319.jpg`

	for i := 0; i < b.N; i++ {
		_, err := gatherFileInfo(d)
		if err != nil {
			b.Fatal(err)
		}

	}
}

func Test_CreateFile(t *testing.T) {
	d := NewDownloader()
	d.fileName = `test1.txt`
	d.fileSize = 23004

	fileInfo, err := setupFile(d)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("file info: %+v", fileInfo)
}

func Benchmark_CreateFile(b *testing.B) {
	b.ReportAllocs()

	d := NewDownloader()
	d.fileName = `test1.txt`
	d.fileSize = 23004

	for i := 0; i < b.N; i++ {
		_, err := setupFile(d)
		if err != nil {
			b.Fatal(err)
		}

	}
}

func Test_CreateChunk(t *testing.T) {
	d := NewDownloader()
	d.fileName = `test1.txt`
	d.fileSize = 23004
	d.numberOfChunks = 5

	err := createChunk(d)
	if err != nil {
		t.Fatal(err)
	}
}

func Benchmark_CreateChunk(b *testing.B) {
	b.ReportAllocs()
	d := NewDownloader()
	d.fileName = `test1.txt`
	d.fileSize = 23004
	d.numberOfChunks = 5

	for i := 0; i < b.N; i++ {
		err := createChunk(d)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Test_download(t *testing.T) {
	err := Download(
		`https://www.pezeshkonline.ir/download/doctormaleki-200519143319.jpg`,
		20,
		0,
		"download")

	if err != nil {
		t.Fatal(err)
	}
}

func Test_CheckMimeType(t *testing.T) {
	CheckMimeType(`jpeg`)
}

func NewDownloader() *Downloader {
	return &Downloader{
		Timeout: 15 * time.Second,
	}
}
