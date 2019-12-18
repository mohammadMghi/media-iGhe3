package handler

import (
	"bufio"
	"fmt"
	"github.com/go-ginger/media/base"
	gm "github.com/go-ginger/models"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	"math"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

type ImageHandler struct {
	DefaultHandler
}

func (h *ImageHandler) GetFilePath(request gm.IRequest) (filePath *base.FilePath, err error) {
	ctx := request.GetContext()
	filePath, err = h.DefaultHandler.GetFilePath(request)
	if err != nil {
		return
	}
	maxStr, exists := ctx.GetQuery("max")
	max, err := strconv.ParseInt(maxStr, 10, 16)
	if err != nil {
		return
	}
	if exists {
		var file *os.File
		file, err = os.Open(filePath.AbsPath)
		if err != nil {
			return
		}
		defer func() {
			err = file.Close()
		}()
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
		newRelativePath := path.Join(
			base.CurrentConfig.CacheDirectoryRelativePath,
			base.CurrentConfig.ImageDirectoryRelativePath,
			filePath.RelativeDirPath,
		)
		filePath.Extension = filepath.Ext(filePath.FullName)
		filePath.Name = filePath.FullName[:len(filePath.FullName)-len(filePath.Extension)]
		newFileFullName := fmt.Sprintf("%s_%dx%d%v", filePath.Name, newWidth, newHeight, filePath.Extension)
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
		var originalImage image.Image
		_, err = file.Seek(0, os.SEEK_SET)
		originalImage, _, err = image.Decode(file)
		if err != nil {
			return
		}
		newImage := resize.Resize(newWidth, newHeight, originalImage, resize.Lanczos3)
		err = h.EnsurePath(newFilePath)
		if err != nil {
			return
		}
		destinationFile, e := os.Create(newFilePath.AbsPath)
		if e != nil {
			err = e
			return
		}
		defer func() {
			err = destinationFile.Close()
		}()
		writer := bufio.NewWriter(destinationFile)
		err = jpeg.Encode(writer, newImage, nil)
		if err != nil {
			return
		}
		err = h.IHandler.SaveFile(file, destinationFile, newFilePath)
		if err != nil {
			return
		}
		filePath = newFilePath
		return
	}
	return
}
