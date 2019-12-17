package upload

import (
	g "github.com/go-ginger/ginger"
)

var CurrentConfig *Config

func Initialize(router *g.RouterGroup, config *Config) {
	CurrentConfig = config
	CurrentConfig.Initialize()
	CurrentConfig.Handler.Initialize(CurrentConfig.Handler)
	upload.Init(upload, nil, nil)
	RegisterRoutes(router)
}
