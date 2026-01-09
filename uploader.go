package main

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/studio-b12/gowebdav"
)

type Uploader struct {
	client    *gowebdav.Client
	sourceDir string
}

func NewUploader(webdavUrl string, webdavUser string, webdavPass string, sourceDir string) (*Uploader, error) {
	client := gowebdav.NewClient(webdavUrl, webdavUser, webdavPass)
	if err := client.Connect(); err != nil {
		return nil, err
	}

	return &Uploader{
		client:    client,
		sourceDir: sourceDir,
	}, nil
}

func (u *Uploader) UploadAll() error {
	// check if anything was downloaded
	if _, err := os.Stat(u.sourceDir); os.IsNotExist(err) {
		return nil
	}

	return filepath.WalkDir(u.sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// get difference between local file path and ./downloads
		remotePath, err := filepath.Rel(u.sourceDir, path)
		if err != nil {
			return err
		}

		if d.IsDir() {
			return u.client.MkdirAll(remotePath, 0o644)
		}
		return u.uploadFile(path, remotePath)
	})
}

func (u *Uploader) uploadFile(localPath string, remotePath string) error {
	// read the file to upload
	file, err := os.Open(localPath)
	if err != nil {
		return err
	}

	defer file.Close()

	return u.client.WriteStream(remotePath, file, 0o644)
}
