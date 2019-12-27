package media

import (
	"github.com/go-m/media/base"
	"github.com/go-m/media/download"
	"github.com/go-m/media/upload"
)

type Config struct {
	base.Config

	Upload   *upload.Config
	Download *download.Config
}
