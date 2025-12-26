package main

type Config struct {
	Playlists []struct {
		Title string `yaml:"title"`
		Link  string `yaml:"link"`
	}
}
