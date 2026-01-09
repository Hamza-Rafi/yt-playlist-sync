package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

type Downloader struct {
	OutputDir     string
	ArchiveFile   string
	BrowserString string
}

func NewDownloader(outputDir string, archiveFile string, browserString string) *Downloader {
	return &Downloader{
		OutputDir:     outputDir,
		ArchiveFile:   archiveFile,
		BrowserString: browserString,
	}
}

func (d *Downloader) DownloadPlaylist(playlist Playlist) error {
	cmd := exec.Command(
		"yt-dlp",
		"--ignore-errors", // yt-dlp returns non 0 exit code on error.
		// one video was private so it exited :c

		"-S", "vcodec:av01,acodec:opus",
		"-f", "bestvideo[height<=480]+bestaudio/best[height<=480]",

		"--merge-output-format", "mp4",
		"-o", filepath.Join(d.OutputDir, playlist.Title, "%(playlist_index)03d - %(title)s.%(ext)s"),

		"--download-archive", d.ArchiveFile,
		"--concurrent-fragments", "2",

		"--embed-thumbnail",
		"--embed-metadata",
		"--embed-subs",
		"--sub-langs", "en,ar",

		"--cookies-from-browser", d.BrowserString,
		playlist.Link,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
