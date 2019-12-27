package media

import (
	g "github.com/go-ginger/ginger"
	"github.com/go-m/media/base"
	"github.com/go-m/media/download"
	"github.com/go-m/media/handler"
	"github.com/go-m/media/upload"
)

type Media struct {
	config *Config

	AuthRouters []*g.RouterGroup
	Router      *g.RouterGroup
	Handlers    map[string]handler.IHandler
}

func (m *Media) Initialize(config *Config, baseConfig interface{}) (err error) {
	base.Initialize(&config.Config)
	if config.Upload == nil {
		config.Upload = new(upload.Config)
	}
	if config.Download == nil {
		config.Download = new(download.Config)
	}
	if m.Handlers == nil {
		imageHandler := handler.ImageHandler{}
		imageHandler.MediaType = &base.MediaType{
			Type:            "image",
			RelativeDirPath: config.ImageDirectoryRelativePath,
		}
		imageHandler.Initialize(&imageHandler)
		m.Handlers = map[string]handler.IHandler{}
		imageHandlerKeys := []string{
			"images",
			"image/jpeg",
			"image/png",
			"image/gif",
			"image/bmp",
			"image/webp",
			"image/vnd.microsoft.icon",
		}
		for _, imageHandlerKey := range imageHandlerKeys {
			m.Handlers[imageHandlerKey] = &imageHandler
		}
	}
	if _, ok := m.Handlers["default"]; !ok {
		defaultHandler := handler.DefaultHandler{}
		defaultHandler.Initialize(&defaultHandler)
		m.Handlers["default"] = &defaultHandler
	}
	handler.CurrentHandlers = m.Handlers
	upload.Initialize(m.Router, config.Upload)
	download.Initialize(m.Router, config.Download)
	return
}
