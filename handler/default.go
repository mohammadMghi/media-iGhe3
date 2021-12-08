package handler

import (
	"crypto/sha256"
	"fmt"
	gm "github.com/go-ginger/models"
	"github.com/go-ginger/models/errors"
	"github.com/mohammadMghi/media-iGhe3/base"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"io"
	"mime/multipart"
	"net/http"
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

func (h *DefaultHandler) SetMediaType(mediaType *base.MediaType) {
	h.MediaType = mediaType
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
		if os.IsNotExist(err) {
			err = errors.GetNotFoundError(request,
				request.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "FileNotFoundError",
						Other: "requested file not found",
					},
				}))
			return
		}
		return
	}
	if _, err = os.Stat(filePath.AbsPath); os.IsNotExist(err) {
		err = errors.GetNotFoundError(request)
		return
	}
	file, err = os.Open(filePath.AbsPath)
	return
}

func (h *DefaultHandler) GetFileName(request gm.IRequest, fileHeader *multipart.FileHeader,
	file multipart.File) (filename string, err error) {
	filename = filepath.Base(fileHeader.Filename)
	return
}

func (h *DefaultHandler) GetPath(request gm.IRequest, fileHeader *multipart.FileHeader,
	file multipart.File) (destinationPath *base.FilePath, err error) {
	mediaType := h.IHandler.GetMediaType()
	hash := sha256.New()
	if _, err = io.Copy(hash, file); err != nil {
		return
	}
	sum := fmt.Sprintf("%x", hash.Sum(nil))
	dirRelativePath := path.Join(mediaType.RelativeDirPath, sum[:2], sum[2:4])
	absDirPath := path.Join(base.CurrentConfig.MediaDirectoryPath, dirRelativePath)
	if _, err = os.Stat(absDirPath); os.IsNotExist(err) {
		err = os.MkdirAll(absDirPath, os.ModePerm)
		if err != nil {
			return
		}
	}
	filename, err := h.IHandler.GetFileName(request, fileHeader, file)
	if err != nil {
		return
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
	destinationPath, err = h.IHandler.GetFilePathWithParams(nil, dirRelativePath, finalFileName)
	if err != nil {
		return
	}
	return
}

func (h *DefaultHandler) Upload(request gm.IRequest, fileHeader *multipart.FileHeader,
	file multipart.File) (fileInfo base.IFileInfo, err error) {
	destinationPath, err := h.IHandler.GetPath(request, fileHeader, file)
	if err != nil {
		return
	}
	reader, _ := file.(io.ReadSeeker)
	fileInfo, err = h.IHandler.SaveFile(reader, nil, destinationPath)
	if err != nil {
		return
	}
	return
}

func (h *DefaultHandler) Download(request gm.IRequest, file *os.File, filePath *base.FilePath) (err error) {
	ctx := request.GetContext()
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	http.ServeContent(ctx.Writer, ctx.Request, filePath.FullName, fileInfo.ModTime(), file)
	return
}
