package media

import (
	g "github.com/go-ginger/ginger"
	"github.com/go-m/media/base"
	"github.com/go-m/media/download"
	"github.com/go-m/media/handler"
	"github.com/go-m/media/upload"
)

type IHandler interface {
	GetBase() (handler *Handler)
	Initialize(handler IHandler, config *Config, handlers map[string]handler.IHandler) (err error)
	InitializeHandlers(config *Config, handlers map[string]handler.IHandler) (err error)
}

type Handler struct {
	IHandler
	config *Config

	AuthRouters []*g.RouterGroup
	Router      *g.RouterGroup
	Handlers    map[string]handler.IHandler
}

func (m *Handler) GetBase() (handler *Handler) {
	return m
}

func (m *Handler) Initialize(handler IHandler, config *Config, handlers map[string]handler.IHandler) (err error) {
	m.IHandler = handler
	base.Initialize(&config.Config)
	if config.Upload == nil {
		config.Upload = new(upload.Config)
	}
	if config.Download == nil {
		config.Download = new(download.Config)
	}
	upload.Initialize(m.Router, config.Upload)
	download.Initialize(m.Router, config.Download)
	err = m.IHandler.InitializeHandlers(config, handlers)
	return
}

func (m *Handler) InitializeHandlers(config *Config, handlers map[string]handler.IHandler) (err error) {
	if m.Handlers == nil {
		m.Handlers = map[string]handler.IHandler{}
	}
	// image handlers
	imageHandler := handler.ImageHandler{}
	imageHandler.Initialize(&imageHandler)
	imageHandler.SetMediaType(&base.MediaType{
		Type:            "image",
		RelativeDirPath: config.ImageDirectoryRelativePath,
	})
	imageHandler.Initialize(&imageHandler)
	imageHandlerKeys := []string{
		config.ImageDirectoryRelativePath,
		"images",
		"image/jpeg",
		"image/png",
		"image/gif",
		"image/bmp",
		"image/webp",
		"image/vnd.microsoft.icon",
	}
	for _, imageHandlerKey := range imageHandlerKeys {
		if _, ok := m.Handlers[imageHandlerKey]; !ok {
			m.Handlers[imageHandlerKey] = &imageHandler
		}
	}
	//
	if handlers != nil {
		for k, h := range handlers {
			m.Handlers[k] = h
		}
	}
	if _, ok := m.Handlers["default"]; !ok {
		defaultHandler := handler.DefaultHandler{}
		defaultHandler.Initialize(&defaultHandler)
		m.Handlers["default"] = &defaultHandler
	}
	handler.CurrentHandlers = m.Handlers
	return
}
