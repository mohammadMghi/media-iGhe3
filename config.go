package media

import (
	"github.com/mohammadMghi/media-iGhe3/base"
	"github.com/mohammadMghi/media-iGhe3/download"
	"github.com/mohammadMghi/media-iGhe3/upload"
)

type Config struct {
	base.Config

	Upload   *upload.Config
	Download *download.Config
}
