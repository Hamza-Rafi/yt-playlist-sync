package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	configFilePath    = flag.String("config", "./config.yml", "path to config file")
	archiveFilePath   = flag.String("archive", "./archive.txt", "path to archive file")
	cookiesFilePath   = flag.String("cookies", "./cookies.txt", "path to cookies file")
	downloadsFilePath = flag.String("downloads", "./downloads/", "downloads folder")
)

func main() {
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Fatalf("loading .env file: %w\n", err)
	}

	config, err := LoadConfig(*configFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	config.CookiesFilePath = *cookiesFilePath
	config.ArchiveFilePath = *archiveFilePath
	config.DownloadsFolderPath = *downloadsFilePath

	runSync(config)
}

func runSync(config Config) {
	// initialise downloader & uploader
	downloader := NewDownloader(config.DownloadsFolderPath, config.ArchiveFilePath, config.CookiesFilePath)
	uploader, err := NewUploader(config.WebDAVURL, config.WebDAVUser, config.WebDAVPass, config.DownloadsFolderPath)
	if err != nil {
		log.Fatalln(err)
	}

	// download playlists
	for i, playlist := range config.Playlists {
		log.Printf("Downloading: %d - %s\n", i, playlist)

		if err := downloader.DownloadPlaylist(playlist); err != nil {
			// if we exit here using log.Fatal, the rest of the playlists won't download
			log.Println(err)
		}
	}

	// upload everything downloaded to the webdav directory
	if err := uploader.UploadAll(); err != nil {
		log.Fatalln(err)
	}

	// delete local files
	os.RemoveAll(config.DownloadsFolderPath)
	log.Println("removed local downloads")
}
