package upload

import (
	gm "github.com/go-ginger/models"
	"github.com/mohammadMghi/media-iGhe3/base"
	"github.com/mohammadMghi/media-iGhe3/handler"
)

type IHandler interface {
	Initialize(handler IHandler)
	Upload(request gm.IRequest) (fileInfo base.IFileInfo, err error)
}

type DefaultHandler struct {
	IHandler
}

func (h *DefaultHandler) Initialize(handler IHandler) {
	h.IHandler = handler
}

func (h *DefaultHandler) Upload(request gm.IRequest) (fileInfo base.IFileInfo, err error) {
	ctx := request.GetContext()
	file, err := ctx.FormFile("file")
	if err != nil {
		return
	}
	f, err := file.Open()
	if err != nil {
		return
	}
	defer func() {
		err = f.Close()
	}()
	fileHandler, err := handler.GetFileHandler(f)
	if err != nil {
		return
	}
	fileInfo, err = fileHandler.Upload(request, file, f)
	return
}
