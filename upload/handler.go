package upload

import (
	"crypto/sha256"
	"fmt"
	"github.com/go-ginger/media/base"
	"github.com/go-ginger/media/handler"
	gm "github.com/go-ginger/models"
	"io"
	"os"
	"path"
	"path/filepath"
)

type IHandler interface {
	Initialize(handler IHandler)
	Upload(request gm.IRequest) (relativePath string, err error)
}

type DefaultHandler struct {
	IHandler
}

func (h *DefaultHandler) Initialize(handler IHandler) {
	h.IHandler = handler
}

func (h *DefaultHandler) Upload(request gm.IRequest) (relativePath string, err error) {
	ctx := request.GetContext()
	file, err := ctx.FormFile("file")
	if err != nil {
		return
	}
	filename := filepath.Base(file.Filename)
	hash := sha256.New()
	f, err := file.Open()
	if err != nil {
		return
	}
	if _, err = io.Copy(hash, f); err != nil {
		return
	}
	sum := fmt.Sprintf("%x", hash.Sum(nil))
	fileHandler, err := handler.GetFileHandler(f)
	if err != nil {
		return
	}
	mediaType := fileHandler.GetMediaType()
	dirRelativePath := path.Join(mediaType.RelativeDirPath, sum[:2], sum[2:4])
	absDirPath := path.Join(base.CurrentConfig.MediaDirectoryPath, dirRelativePath)
	if _, err = os.Stat(absDirPath); os.IsNotExist(err) {
		err = os.MkdirAll(absDirPath, os.ModePerm)
		if err != nil {
			return
		}
	}
	finalFileName := filename
	var finalFilePath string
	fileNumber := 1
	for {
		finalFilePath = path.Join(absDirPath, finalFileName)
		if _, err := os.Stat(finalFilePath); os.IsNotExist(err) {
			break
		}
		fileNumber++
		ext := filepath.Ext(filename)
		name := filename[:len(filename)-len(ext)]
		if ext != "" {
			finalFileName = fmt.Sprintf("%v_%v%v", name, fileNumber, ext)
		} else {
			finalFileName = fmt.Sprintf("%v_%v", name, fileNumber)
		}
	}
	err = ctx.SaveUploadedFile(file, finalFilePath)
	if err != nil {
		return
	}
	relativePath = path.Join(dirRelativePath, finalFileName)
	return
}
