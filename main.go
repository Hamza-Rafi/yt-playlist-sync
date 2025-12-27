package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

var (
	configFilePath      string = "./config.yml"
	archiveFilePath     string = "./archive.txt"
	downloadsFolderPath string = "./downloads/"
)

func main() {
	// load config from config file
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalln(err)
	}

	for i, playlist := range config.Playlists {
		fmt.Print("Downloading: ")
		fmt.Printf("%d - %s - %s\n", i, playlist.Title, playlist.Link)

		cmd := exec.Command(
			"yt-dlp",

			"-S", "vcodec:av01,acodec:opus",
			"-f", "bestvideo[height<=480]+bestaudio/best[height<=480]",

			"--merge-output-format", "mp4",
			"-o", filepath.Join(downloadsFolderPath, playlist.Title, "%(playlist_index)03d - %(title)s.%(ext)s"),

			"--download-archive", archiveFilePath,
			"--concurrent-fragments", "4",

			"--embed-thumbnail",
			"--embed-metadata",
			"--restrict-filenames",

			"--cookies-from-browser", "firefox:/home/hamza/.zen/449e99hw.Default (beta)",
			playlist.Link,
		)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Fatalln(err)
		}
	}
}
