package upload

import (
	"github.com/go-ginger/ginger"
	gm "github.com/go-ginger/models"
)

type uploadController struct {
	ginger.BaseItemsController
}

func (c *uploadController) Post(request gm.IRequest) (result interface{}) {
	ctx := request.GetContext()
	fileInfo, err := CurrentConfig.Handler.Upload(request)
	if c.HandleErrorNoResult(request, err) {
		return
	}
	ctx.JSON(201, fileInfo)
	return
}
