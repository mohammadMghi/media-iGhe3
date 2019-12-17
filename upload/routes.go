package upload

import (
	g "github.com/go-ginger/ginger"
)

var upload = new(uploadController)

func RegisterRoutes(router *g.RouterGroup) {
	handlers := make([]g.HandlerFunc, 0)
	if CurrentConfig.MustHaveRoles != nil {
		handlers = append(handlers, CurrentConfig.LoginHandler.MustHaveRole(CurrentConfig.MustHaveRoles...))
	} else if CurrentConfig.MustAuthenticate {
		handlers = append(handlers, CurrentConfig.LoginHandler.MustAuthenticate())
	}
	upload.AddRoute("Post", handlers...)
	upload.RegisterRoutes(upload, "/upload", router)
}
