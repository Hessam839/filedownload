package filedownload

import (
	"testing"
	"time"
)

func Test_GetFileInfo(t *testing.T) {
	d := NewDownloader()
	d.uri = `https://dl2.soft98.ir/soft/p-q/PassMark.BurnInTest.Professional.9.2.Build.1007.rar?1623043852`
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
	var values []int64 = []int64{23004, 13_000_000, 5_604_002}
	d := NewDownloader()
	d.fileName = `test1.txt`

	for _, v := range values {
		d.fileSize = v
		err := createChunk(d)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("for size: %d number of chunk is: %d", d.fileSize, d.numberOfChunks)
	}
}

func Benchmark_CreateChunk(b *testing.B) {
	b.ReportAllocs()
	d := NewDownloader()
	d.fileName = `test1.txt`
	d.fileSize = 23004

	for i := 0; i < b.N; i++ {
		err := createChunk(d)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Test_download(t *testing.T) {
	err := Download(
		`https://dl2.soft98.ir/soft/n/Notepad.8.0.x86.rar?1623045509`,
		6,
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
		Timeout:             15 * time.Second,
		numberOfChunks:      100,
		numberOfConnections: 20,
	}
}
