package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
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
	uploader, err := NewUploader(config.WebDAVURL, config.WebDAVUser, config.WebDAVPass, downloadsFolderPath)
	if err != nil {
		log.Fatalln(err)
	}

	// download playlists
	for i, playlist := range config.Playlists {
		log.Printf("Downloading: %d - %s\n", i, playlist.Title)

		if err := downloader.DownloadPlaylist(playlist); err != nil {
			// if we exit here using log.Fatal, the rest of the playlists won't download
			log.Println(err)
		}
	}

	// upload everything to the webdav directory
	if err := uploader.UploadAll(); err != nil {
		log.Fatalln(err)
	}

	// delete local files
	os.RemoveAll(downloadsFolderPath)
	log.Println("removed local downloads")
}
