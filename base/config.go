package base

import (
	"os"
	"path"
)

type Config struct {
	MediaDirectoryPath         string
	CacheDirectoryRelativePath string
	CacheDirectoryPath         string
	FileDirectoryRelativePath  string
	FileDirectoryPath          string
	ImageDirectoryRelativePath string
	ImageDirectoryPath         string
}

var CurrentConfig *Config

func (c *Config) Initialize() {
	CurrentConfig = c
	if c.MediaDirectoryPath == "" {
		c.MediaDirectoryPath = path.Join(os.Getenv("HOME"), "media")
	}
	if c.CacheDirectoryRelativePath == "" {
		c.CacheDirectoryRelativePath = ".cache"
	}
	if c.ImageDirectoryRelativePath == "" {
		c.ImageDirectoryRelativePath = "images"
	}
	if c.FileDirectoryRelativePath == "" {
		c.FileDirectoryRelativePath = "files"
	}
	c.ImageDirectoryPath = path.Join(c.MediaDirectoryPath, c.ImageDirectoryRelativePath)
	c.FileDirectoryPath = path.Join(c.MediaDirectoryPath, c.FileDirectoryRelativePath)
	c.CacheDirectoryPath = path.Join(c.MediaDirectoryPath, c.CacheDirectoryRelativePath)
}
