package main

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	WebDAVURL  string `yaml:"webdav_url"`
	WebDAVUser string
	WebDAVPass string

	BrowserString string `yaml:"browser_string"`

	Playlists []Playlist
}

type Playlist struct {
	Title string `yaml:"title"`
	Link  string `yaml:"link"`
}

func LoadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("reading config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("error unmarshalling config: %w", err)
	}

	config.WebDAVUser = os.Getenv("WEBDAV_USER")
	config.WebDAVPass = os.Getenv("WEBDAV_PASS")

	if config.WebDAVUser == "" || config.WebDAVPass == "" {
		return Config{}, fmt.Errorf("Ensure WEBDAV_USER and WEBDAV_PASS env vars are set")
	}

	return config, nil
}
