package handler

import (
	"github.com/go-ginger/media/base"
	gm "github.com/go-ginger/models"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type IHandler interface {
	Initialize(handler IHandler)
	GetMediaType() (mediaType *base.MediaType)
	GetFilePath(request gm.IRequest) (filePath *base.FilePath, err error)
	GetFilePathWithParams(mediaType *string, relativeDirPath, fileName string) (filePath *base.FilePath, err error)
	GetFile(request gm.IRequest) (fileName string, file *os.File, err error)
	SaveFile(file, destinationFile *os.File, destination *base.FilePath) (err error)
}

var CurrentHandlers map[string]IHandler

func GetHandlerByKey(key string) (handler IHandler, ok bool) {
	handler, ok = CurrentHandlers[key]
	return
}

func GetFileHandler(file multipart.File) (fileHandler IHandler, err error) {
	buffer := make([]byte, 512)
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return
	}
	_, err = file.Read(buffer)
	if err != nil {
		return
	}
	contentType := http.DetectContentType(buffer)
	fileHandler, ok := GetHandlerByKey(contentType)
	if !ok {
		fileHandler, _ = GetHandlerByKey("default")
	}
	return
}
