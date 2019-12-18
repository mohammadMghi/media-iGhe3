package download

import (
	"fmt"
	mb "github.com/go-ginger/media/base"
	"github.com/go-ginger/media/handler"
	gm "github.com/go-ginger/models"
	"github.com/go-ginger/models/errors"
	"io"
	"net/http"
	"os"
)

type IHandler interface {
	Initialize(handler IHandler)
	Download(request gm.IRequest) (err error)
	GetFile(request gm.IRequest) (filePath *mb.FilePath, file *os.File, err error)
}

type DefaultHandler struct {
	IHandler
}

func (h *DefaultHandler) Initialize(iHandler IHandler) {
	h.IHandler = iHandler
}

func (h *DefaultHandler) GetFile(request gm.IRequest) (filePath *mb.FilePath, file *os.File, err error) {
	req := request.GetBaseRequest()
	mediaType, _ := req.Params.Get("media_type")
	currentHandler, ok := handler.GetHandlerByKey(mediaType)
	if !ok {
		err = errors.GetValidationError("media handler for this media type not found")
		return
	}
	return currentHandler.GetFile(request)
}

func (h *DefaultHandler) Download(request gm.IRequest) (err error) {
	filePath, file, err := h.GetFile(request)
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
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	http.ServeContent(ctx.Writer, ctx.Request, filePath.FullName, fileInfo.ModTime(), file)
	return
}
