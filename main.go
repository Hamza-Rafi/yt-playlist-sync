package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/studio-b12/gowebdav"
)

var (
	configFilePath      string = "./config.yml"
	archiveFilePath     string = "./archive.txt"
	downloadsFolderPath string = "./downloads/"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("loading .env file: %w", err)
	}

	// load config from config file
	config, err := LoadConfig(configFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	downloader := NewDownloader(downloadsFolderPath, archiveFilePath, config.BrowserString)

	// download playlists
	for i, playlist := range config.Playlists {
		log.Printf("Downloading: %d - %s\n", i, playlist.Title)

		if err := downloader.DownloadPlaylist(playlist); err != nil {
			// if we exit here using log.Fatal, the rest of the playlists won't download
			log.Println(err)
		}
	}

	webdavClient := gowebdav.NewClient(config.WebDAVURL, config.WebDAVUser, config.WebDAVPass)
	if err := webdavClient.Connect(); err != nil {
		log.Fatalln("error connecting to webdav server: ", err)
	}

	filepath.WalkDir(downloadsFolderPath, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// get difference between filepath and ./downloads
		remotePath, err := filepath.Rel(downloadsFolderPath, filePath)

		if d.IsDir() {
			err := webdavClient.MkdirAll(remotePath, 0o644)

			log.Println("mkdir: ", remotePath)
			if err != nil {
				log.Fatalln("error creating: ", remotePath, err)
			}
		} else {
			err := uploadFile(webdavClient, filePath, remotePath)
			if err != nil {
				log.Println("error writing: ", remotePath, err)
			} else {
				log.Println("uploaded: ", remotePath)
			}
		}

		return nil
	})

	// delete local files
	os.RemoveAll(downloadsFolderPath)
	log.Println("removed local downloads")
}

func uploadFile(webdavClient *gowebdav.Client, filePath string, remotePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	return webdavClient.WriteStream(remotePath, file, 0o644)
}
