package handler

import (
	gm "github.com/go-ginger/models"
	"github.com/go-m/media/base"
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
	GetFile(request gm.IRequest) (filePath *base.FilePath, file *os.File, err error)
	SaveFile(file io.ReadSeeker, destinationFile *os.File, destination *base.FilePath) (fileInfo base.IFileInfo, err error)
}

var CurrentHandlers map[string]IHandler

func GetHandlerByKey(key string) (handler IHandler) {
	handler, ok := CurrentHandlers[key]
	if !ok {
		handler, _ = CurrentHandlers["default"]
	}
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
	fileHandler = GetHandlerByKey(contentType)
	return
}
