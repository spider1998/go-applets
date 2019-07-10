package routing

import (
	"fmt"
	"net/http"
	"sdkeji/wechat/pkg/app"
	open "sdkeji/wechat/pkg/routing/app"
	"sdkeji/wechat/pkg/routing/internal"

	"git.sdkeji.top/share/sdlib/log"

	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/content"
)

func Register(logger log.Logger) http.Handler {
	router := routing.New()
	router.NotFound(notFound)
	router.Use(
		routingLogger(logger),
		errorHandler(logger),
		content.TypeNegotiator(content.JSON),
	)

	api := router.Group("/" + app.System)

	{
		versionHandler := NewVersionHandler()
		api.Get("/version", versionHandler.Version)
		open.Register(api.Group("/app/v1"))
		internal.Register(api.Group("/internal/v1"))

	}

	for _, route := range router.Routes() {
		logger.Debug(fmt.Sprintf("register route: \"%-6s -> %s\".", route.Method(), route.Path()))
	}

	return router
}
