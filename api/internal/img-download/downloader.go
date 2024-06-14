package img_download

import (
	"io"
	"net/http"
	"os"
)

type ImgDownloader struct {
	downloadDir string
}

var imgDownloader ImgDownloader

func GetDownloader() *ImgDownloader {
	return &imgDownloader
}

func (*ImgDownloader) InitializeDownloader(downloadDir string) error {
	imgDownloader = ImgDownloader{downloadDir}
	return nil
}

func Download(downloader *ImgDownloader, url string, fileName string) (string, error) {
	filePath := downloader.downloadDir + "/" + fileName
	create, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer create.Close()

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	_, err = io.Copy(create, resp.Body)
	if err != nil {
		return "", err
	}
	return filePath, nil
}
