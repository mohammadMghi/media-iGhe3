package download

import (
	g "github.com/go-ginger/ginger"
)

var CurrentConfig *Config

func Initialize(controller g.IController, router *g.RouterGroup, config *Config) {
	CurrentConfig = config
	CurrentConfig.Initialize()
	CurrentConfig.Handler.Initialize(CurrentConfig.Handler)
	if controller == nil {
		controller = download
	}
	download.Init(controller, nil, nil)
	RegisterRoutes(router)
}
