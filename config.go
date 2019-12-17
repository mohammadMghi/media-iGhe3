package media

import (
	g "github.com/go-ginger/ginger"
	"github.com/go-ginger/media/base"
	"github.com/go-ginger/media/download"
	"github.com/go-ginger/media/upload"
)

type Config struct {
	base.Config

	AuthRouters []*g.RouterGroup
	Router      *g.RouterGroup
	Upload      *upload.Config
	Download    *download.Config
}
