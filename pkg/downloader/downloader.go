/*
Package downloader ...
Package to download a given file from URL and return its size after download
*/
package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

//GetFile ...
//Download file and place it in the destination path with custom name
func GetFile(file, dest, url string) (int64, error) {
	fmt.Printf("* Downloading %s...", file)
	size, err := GetFileSize(url)
	if err != nil {
		return 0, err
	}
	res, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	dest = filepath.Join(dest, file)
	fp, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0755)
	if err != nil {
		res.Body.Close()
		return 0, err
	}
	written, err := io.Copy(fp, res.Body)
	if err != nil {
		fp.Close()
		res.Body.Close()
		return 0, err
	}
	fp.Sync()
	fp.Close()
	res.Body.Close()
	fileInfo, err := os.Stat(dest)
	fmt.Println("Done , size - ", fileInfo.Size(), "Bytes")
	if size != written {
		return written, err
	}
	return written, nil
}

//GetFileSize ...
//Returns the size of a file in bytes
func GetFileSize(url string) (int64, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return 0, err
	}
	size := resp.Header.Get("Content-Length")
	sizeInt, err := strconv.Atoi(size)
	if err != nil {
		return 0, err
	}
	return int64(sizeInt), nil

}
