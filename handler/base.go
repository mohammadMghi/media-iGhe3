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
	SetMediaType(mediaType *base.MediaType)
	GetMediaType() (mediaType *base.MediaType)
	GetFilePath(request gm.IRequest) (filePath *base.FilePath, err error)
	GetFilePathWithParams(mediaType *string, relativeDirPath, fileName string) (filePath *base.FilePath, err error)
	GetFile(request gm.IRequest) (filePath *base.FilePath, file *os.File, err error)
	SaveFile(file io.ReadSeeker, destinationFile *os.File, destination *base.FilePath) (fileInfo base.IFileInfo, err error)
	GetFileName(request gm.IRequest, fileHeader *multipart.FileHeader, file multipart.File) (filename string, err error)
	GetPath(request gm.IRequest, fileHeader *multipart.FileHeader, file multipart.File) (destinationPath *base.FilePath, err error)
	Upload(request gm.IRequest, fileHeader *multipart.FileHeader, file multipart.File) (fileInfo base.IFileInfo, err error)
	Download(request gm.IRequest, file *os.File, filePath *base.FilePath) (err error)
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
