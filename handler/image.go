package handler

import (
	"fmt"
	"github.com/go-ginger/media/base"
	gm "github.com/go-ginger/models"
	"github.com/nfnt/resize"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

type ImageHandler struct {
	DefaultHandler
}

func (h *ImageHandler) EnsureImageMaxSize(file io.ReadSeeker, max int64) (out io.ReadSeeker,
	width, height uint, err error) {
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return
	}
	if max <= 0 {
		out = file
		return
	}
	var imageConf image.Config
	imageConf, _, err = image.DecodeConfig(file)
	if err != nil {
		return
	}
	resizeRatio := math.Max(
		math.Min(float64(max), float64(imageConf.Width))/float64(imageConf.Width),
		math.Min(float64(max), float64(imageConf.Height))/float64(imageConf.Height),
	)
	if resizeRatio >= 1 {
		out = file
		return
	}
	width = uint(float64(imageConf.Width) * resizeRatio)
	height = uint(float64(imageConf.Height) * resizeRatio)
	var originalImage image.Image
	_, err = file.Seek(0, io.SeekStart)
	originalImage, _, err = image.Decode(file)
	if err != nil {
		return
	}
	newImage := resize.Resize(width, height, originalImage, resize.Lanczos3)
	var rws base.ReadWriterSeeker
	rws.InitializeWriter()
	err = jpeg.Encode(&rws, newImage, nil)
	if err != nil {
		return
	}
	out = rws.GetReadSeeker()
	_, err = out.Seek(0, io.SeekStart)
	return
}

func (h *ImageHandler) SaveFile(file io.ReadSeeker, destinationFile *os.File,
	destination *base.FilePath) (fileInfo base.IFileInfo, err error) {
	fileInfo, err = h.DefaultHandler.SaveFile(file, destinationFile, destination)
	if err != nil {
		return
	}
	var savedFile *os.File
	savedFile, err = os.Open(destination.AbsPath)
	if err != nil {
		return
	}
	defer func() {
		err = savedFile.Close()
	}()
	var imageConf image.Config
	imageConf, _, err = image.DecodeConfig(savedFile)
	if err != nil {
		return
	}
	uploadedFileInfo := fileInfo.(*base.FileInfo)
	uploadedFileInfo.Width = &imageConf.Width
	uploadedFileInfo.Height = &imageConf.Height
	return
}

func (h *ImageHandler) GetFilePath(request gm.IRequest) (filePath *base.FilePath, err error) {
	ctx := request.GetContext()
	filePath, err = h.IHandler.GetFilePath(request)
	if err != nil {
		return
	}
	maxStr, exists := ctx.GetQuery("max")
	if exists {
		var max int64
		max, err = strconv.ParseInt(maxStr, 10, 16)
		if err != nil {
			return
		}
		var file *os.File
		file, err = os.Open(filePath.AbsPath)
		if err != nil {
			return
		}
		defer func() {
			err = file.Close()
		}()
		newRelativePath := path.Join(
			base.CurrentConfig.CacheDirectoryRelativePath,
			base.CurrentConfig.ImageDirectoryRelativePath,
			filePath.RelativeDirPath,
		)
		filePath.Extension = filepath.Ext(filePath.FullName)
		filePath.Name = filePath.FullName[:len(filePath.FullName)-len(filePath.Extension)]

		var imageConf image.Config
		imageConf, _, err = image.DecodeConfig(file)
		if err != nil {
			return
		}
		resizeRatio := math.Max(
			math.Min(float64(max), float64(imageConf.Width))/float64(imageConf.Width),
			math.Min(float64(max), float64(imageConf.Height))/float64(imageConf.Height),
		)
		if resizeRatio >= 1 {
			return
		}
		newWidth := uint(float64(imageConf.Width) * resizeRatio)
		newHeight := uint(float64(imageConf.Height) * resizeRatio)

		newFileFullName := fmt.Sprintf("%s_%dx%d%s", filePath.Name, newWidth, newHeight, filePath.Extension)
		var newFilePath *base.FilePath
		newFilePath, err = h.IHandler.GetFilePathWithParams(nil, newRelativePath, newFileFullName)
		if err != nil {
			return
		}
		if _, e := os.Stat(newFilePath.AbsPath); e == nil || os.IsExist(e) {
			// already exists
			filePath = newFilePath
			return
		}
		reader, newWidth, newHeight, e := h.EnsureImageMaxSize(file, max)
		if e != nil {
			err = e
			return
		}
		if newWidth <= 0 {
			filePath = newFilePath
			return
		}
		_, err = h.IHandler.SaveFile(reader, nil, newFilePath)
		if err != nil {
			return
		}
		filePath = newFilePath
		return
	}
	return
}
