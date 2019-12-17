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
	downloadRelativeURL, err := CurrentConfig.Handler.Upload(request)
	if c.HandleErrorNoResult(request, err) {
		return
	}
	ctx.String(200, downloadRelativeURL)
	return
}
