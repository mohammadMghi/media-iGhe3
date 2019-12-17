package base

import (
	"os"
	"path"
)

type Config struct {
	MediaDirectoryPath         string
	ImageDirectoryRelativePath string
	ImageDirectoryPath         string
}

var CurrentConfig *Config

func (c *Config) Initialize() {
	CurrentConfig = c
	if c.MediaDirectoryPath == "" {
		c.MediaDirectoryPath = path.Join(os.Getenv("HOME"), "media")
	}
	if c.ImageDirectoryRelativePath == "" {
		c.ImageDirectoryRelativePath = "images"
	}
	c.ImageDirectoryPath = path.Join(c.MediaDirectoryPath, c.ImageDirectoryRelativePath)
}
