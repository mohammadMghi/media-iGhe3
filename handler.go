package media

import (
	"github.com/go-ginger/media/base"
	"github.com/go-ginger/media/download"
	"github.com/go-ginger/media/upload"
)

type Handler struct {
	config *Config
}

func (h *Handler) Initialize(config *Config, baseConfig interface{}) (err error) {
	base.Initialize(&config.Config)
	if config.Upload == nil {
		config.Upload = new(upload.Config)
	}
	if config.Download == nil {
		config.Download = new(download.Config)
	}
	upload.Initialize(config.Router, config.Upload)
	download.Initialize(config.Router, config.Download)
	return
}
