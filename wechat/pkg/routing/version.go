package routing

import (
	"github.com/go-ozzo/ozzo-routing"
	"sdkeji/wechat/pkg/app"
)

func NewVersionHandler() VersionHandler {
	return VersionHandler{}
}

type VersionHandler struct{}

func (h VersionHandler) Version(c *routing.Context) error {
	return c.Write(map[string]string{
		"version":    app.Version,
		"build_time": app.BuildTime,
	})
}
