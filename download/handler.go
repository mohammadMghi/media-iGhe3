package download

import (
	"fmt"
	gm "github.com/go-ginger/models"
	"github.com/go-m/media/handler"
	"io"
)

type IHandler interface {
	Initialize(handler IHandler)
	Download(request gm.IRequest) (err error)
}

type DefaultHandler struct {
	IHandler
}

func (h *DefaultHandler) Initialize(iHandler IHandler) {
	h.IHandler = iHandler
}

func (h *DefaultHandler) Download(request gm.IRequest) (err error) {
	req := request.GetBaseRequest()
	mediaType, _ := req.Params.Get("media_type")
	currentHandler := handler.GetHandlerByKey(mediaType)
	filePath, file, err := currentHandler.GetFile(request)
	if err != nil {
		return
	}
	defer func() {
		err = file.Close()
	}()
	ctx := request.GetContext()
	if CurrentConfig.DownloadAsAttachment {
		ctx.Writer.Header().Set("Content-Disposition",
			fmt.Sprintf("attachment; filename=%v", filePath.FullName))
		ctx.Writer.Header().Set("Content-Type", ctx.Request.Header.Get("Content-Type"))
		_, err = io.Copy(ctx.Writer, file)
		if err != nil {
			return
		}
		return
	}
	err = currentHandler.Download(request, file, filePath)
	return
}
