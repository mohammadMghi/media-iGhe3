package upload

import ab "github.com/go-m/auth/base"

type Config struct {
	LoginHandler ab.ILoginHandler
	Handler      IHandler

	MustAuthenticate bool
	MustHaveRoles    []string
}

func (c *Config) Initialize() {
	if c.Handler == nil {
		c.Handler = &DefaultHandler{}
	}
}
