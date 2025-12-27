package main

type Config struct {
	WebDAVURL  string `yaml:"webdav_url"`
	WebDAVUser string
	WebDAVPass string

	Playlists []Playlist
}

type Playlist struct {
	Title string `yaml:"title"`
	Link  string `yaml:"link"`
}
