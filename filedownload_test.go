package filedownload

import "testing"

func Test_GetFileInfo(t *testing.T) {
	d := NewDownloader()
	d.uri = `https://www.pezeshkonline.ir/download/doctormaleki-200519143319.jpg`
	info, err := gatherFileInfo(d)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("url info: %+v", info)
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

func Test_download(t *testing.T) {
	err := Download(`https://www.pezeshkonline.ir/download/doctormaleki-200519143319.jpg`, 10)
	if err != nil {
		t.Fatal(err)
	}

}