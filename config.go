package media

import (
	"github.com/go-ginger/media/base"
	"github.com/go-ginger/media/download"
	"github.com/go-ginger/media/upload"
)

type Config struct {
	base.Config

	Upload   *upload.Config
	Download *download.Config
}
