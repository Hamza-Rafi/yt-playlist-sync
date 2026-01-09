package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
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

	// download playlists
	for i, playlist := range config.Playlists {
		fmt.Print("Downloading: ")
		fmt.Printf("%d - %s - %s\n", i, playlist.Title, playlist.Link)

		err := downloadPlaylist(playlist)
		if err != nil {
			log.Println("error downloading playlist: ", err)
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

func downloadPlaylist(playlist Playlist) error {
	cmd := exec.Command(
		"yt-dlp",
		"--ignore-errors", // yt-dlp returns non 0 exit code on error.
		// one video was private so it exited :c

		"-S", "vcodec:av01,acodec:opus",
		"-f", "bestvideo[height<=480]+bestaudio/best[height<=480]",

		"--merge-output-format", "mp4",
		"-o", filepath.Join(downloadsFolderPath, playlist.Title, "%(playlist_index)03d - %(title)s.%(ext)s"),

		"--download-archive", archiveFilePath,
		"--concurrent-fragments", "2",

		"--embed-thumbnail",
		"--embed-metadata",
		"--embed-subs",
		"--sub-langs", "en,ar",

		"--cookies-from-browser", "firefox:/home/hamza/.zen/449e99hw.Default (beta)",
		playlist.Link,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func uploadFile(webdavClient *gowebdav.Client, filePath string, remotePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	return webdavClient.WriteStream(remotePath, file, 0o644)
}
