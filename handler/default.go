package handler

import (
	"github.com/go-m/media/base"
	gm "github.com/go-ginger/models"
	"github.com/go-ginger/models/errors"
	"io"
	"os"
	"path"
	"path/filepath"
)

type DefaultHandler struct {
	IHandler
	MediaType *base.MediaType
}

func (h *DefaultHandler) Initialize(handler IHandler) {
	h.IHandler = handler
	if h.MediaType == nil {
		h.MediaType = &base.MediaType{
			Type:            "default",
			RelativeDirPath: "files",
		}
	}
}

func (h *DefaultHandler) EnsurePath(path *base.FilePath) (err error) {
	if _, e := os.Stat(path.AbsDirPath); os.IsNotExist(e) {
		err = os.MkdirAll(path.AbsDirPath, os.ModePerm)
		if err != nil {
			return
		}
	}
	return
}

func (h *DefaultHandler) SaveFile(reader io.ReadSeeker, destinationFile *os.File,
	destination *base.FilePath) (fileInfo base.IFileInfo, err error) {
	_, err = reader.Seek(0, io.SeekStart)
	if err != nil {
		return
	}
	err = h.EnsurePath(destination)
	if err != nil {
		return
	}
	if destinationFile == nil {
		destinationFile, err = os.Create(destination.AbsPath)
		if err != nil {
			return
		}
		defer func() {
			e := destinationFile.Close()
			if e != nil {
				err = e
			}
		}()
	}
	_, err = io.Copy(destinationFile, reader)
	fileInfo = &base.FileInfo{
		RelativeURL: destination.RelativePath,
	}
	return
}

func (h *DefaultHandler) GetMediaType() (mediaType *base.MediaType) {
	mediaType = h.MediaType
	return
}

func (h *DefaultHandler) GetFilePathWithParams(mediaType *string, relativeDirPath, fileName string) (filePath *base.FilePath,
	err error) {
	filePath = &base.FilePath{}
	filePath.RelativePath = path.Join(relativeDirPath, fileName)
	filePath.FullName = fileName
	filePath.RelativeDirPath = relativeDirPath
	filePath.AbsDirPath = base.CurrentConfig.MediaDirectoryPath
	if mediaType != nil {
		filePath.AbsDirPath = path.Join(filePath.AbsDirPath, *mediaType)
	}
	filePath.AbsDirPath = path.Join(filePath.AbsDirPath, filePath.RelativeDirPath)
	filePath.AbsPath = path.Join(filePath.AbsDirPath, filePath.FullName)
	filePath.Extension = filepath.Ext(filePath.FullName)
	filePath.Name = filePath.FullName[:len(filePath.FullName)-len(filePath.Extension)]
	return
}

func (h *DefaultHandler) GetFilePath(request gm.IRequest) (filePath *base.FilePath, err error) {
	req := request.GetBaseRequest()
	mediaType, _ := req.Params.Get("media_type")
	var relativeDirPath string
	var fileFullName string
	for _, param := range []string{"p1", "p2", "p3", "p4", "p5"} {
		if paramValue, exists := req.Params.Get(param); exists {
			relativeDirPath = path.Join(relativeDirPath, paramValue)
			fileFullName = paramValue
		} else {
			break
		}
	}
	relativeDirPath = path.Dir(relativeDirPath)
	filePath, err = h.IHandler.GetFilePathWithParams(&mediaType, relativeDirPath, fileFullName)
	return
}

func (h *DefaultHandler) GetFile(request gm.IRequest) (filePath *base.FilePath, file *os.File, err error) {
	filePath, err = h.IHandler.GetFilePath(request)
	if err != nil {
		return
	}
	if _, err = os.Stat(filePath.AbsPath); os.IsNotExist(err) {
		err = errors.GetNotFoundError()
		return
	}
	file, err = os.Open(filePath.AbsPath)
	return
}
