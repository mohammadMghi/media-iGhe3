package download

import (
	"github.com/go-ginger/ginger"
	gm "github.com/go-ginger/models"
)

type downloadController struct {
	ginger.BaseItemsController
}

func (c *downloadController) Get(request gm.IRequest) (result interface{}) {
	ctx := request.GetContext()
	err := CurrentConfig.Handler.Download(request)
	if c.HandleErrorNoResult(request, err) {
		return
	}
	ctx.Status(204)
	return
}
