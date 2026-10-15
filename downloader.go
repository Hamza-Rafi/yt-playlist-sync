package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Downloader struct {
	OutputDir   string
	ArchiveFile string
	CookiesFile string
}

func NewDownloader(outputDir string, archiveFile string, cookiesFile string) *Downloader {
	return &Downloader{
		OutputDir:   outputDir,
		ArchiveFile: archiveFile,
		CookiesFile: cookiesFile,
	}
}

func (d *Downloader) DownloadPlaylist(playlistLink string) error {
	fmt.Println(d.CookiesFile)
	cmd := exec.Command(
		"yt-dlp",
		"--ignore-errors", // yt-dlp returns non 0 exit code on error.
		// one video was private so it exited :c
		"--continue",

		"-S", "vcodec:av01,acodec:opus",
		"-f", "bestvideo[height<=480]+bestaudio/best[height<=480]",

		"--merge-output-format", "mp4",
		"-o", filepath.Join(d.OutputDir, "%(playlist_title)s/%(playlist_index)03d - %(title)s.%(ext)s"),

		"--download-archive", d.ArchiveFile,
		"--concurrent-fragments", "2",

		"--write-subs",
		"--sub-langs", "en,ar",
		"--embed-subs",

		"--embed-thumbnail",
		"--embed-metadata",

		"--cookies", d.CookiesFile,
		playlistLink,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
