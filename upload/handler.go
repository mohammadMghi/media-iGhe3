package upload

import (
	"crypto/sha256"
	"fmt"
	"github.com/go-ginger/media/base"
	gm "github.com/go-ginger/models"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
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
	dirRelativePath := path.Join(base.CurrentConfig.ImageDirectoryRelativePath, sum[:2], sum[2:4])
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
		name := filename
		parts := strings.Split(filename, ".")
		if len(parts) > 0 {
			extension := parts[len(parts)-1]
			name = filename[:len(filename)-len(extension)-1]
			finalFileName = fmt.Sprintf("%v_%v.%v", name, fileNumber, extension)
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
