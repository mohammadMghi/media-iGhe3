package download

import (
	g "github.com/go-ginger/ginger"
)

var download = new(downloadController)

func RegisterRoutes(router *g.RouterGroup) {
	handlers := make([]g.HandlerFunc, 0)
	if CurrentConfig.MustHaveRoles != nil {
		handlers = append(handlers, CurrentConfig.LoginHandler.MustHaveRole(CurrentConfig.MustHaveRoles...))
	} else if CurrentConfig.MustAuthenticate {
		handlers = append(handlers, CurrentConfig.LoginHandler.MustAuthenticate())
	}
	download.AddRoute("Get", handlers...)
	download.RegisterRoutes(download, "/download/:media_type/:p1", router)
	download.RegisterRoutes(download, "/download/:media_type/:p1/:p2", router)
	download.RegisterRoutes(download, "/download/:media_type/:p1/:p2/:p3", router)
	download.RegisterRoutes(download, "/download/:media_type/:p1/:p2/:p3/:p4", router)
	download.RegisterRoutes(download, "/download/:media_type/:p1/:p2/:p3/:p4/:p5", router)
}
