package download

import (
	"fmt"
	"github.com/go-ginger/media/base"
	gm "github.com/go-ginger/models"
	"io"
	"net/http"
	"os"
	"path"
)

type IHandler interface {
	Initialize(handler IHandler)
	Download(request gm.IRequest) (err error)
}

type DefaultHandler struct {
	IHandler
}

func (h *DefaultHandler) Initialize(handler IHandler) {
	h.IHandler = handler
}

func (h *DefaultHandler) Download(request gm.IRequest) (err error) {
	req := request.GetBaseRequest()
	relativePath := ""
	fileName := ""
	for _, param := range []string{"media_type", "p1", "p2", "p3", "p4", "p5"} {
		if paramValue, exists := req.Params.Get(param); exists {
			relativePath = path.Join(relativePath, paramValue)
			fileName = paramValue
		} else {
			break
		}
	}
	fullPath := path.Join(base.CurrentConfig.MediaDirectoryPath, relativePath)
	ctx := request.GetContext()
	file, err := os.Open(fullPath)
	if err != nil {
		return
	}
	defer func() {
		err = file.Close()
	}()
	if CurrentConfig.DownloadAsAttachment {
		ctx.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%v", fileName))
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
	http.ServeContent(ctx.Writer, ctx.Request, fileName, fileInfo.ModTime(), file)
	return
}
