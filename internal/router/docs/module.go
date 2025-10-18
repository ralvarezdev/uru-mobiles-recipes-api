package router

import (
	gonethttp "github.com/ralvarezdev/go-net/http"
	"github.com/swaggo/http-swagger"
)

var (
	Module = &gonethttp.Module{
		Pattern: "/docs",
		AddHandlersFn: func(m *gonethttp.Module) {
			m.AddHandleFunc(
				"GET /swagger/*",
				httpSwagger.WrapHandler,
			)
		},
	}
)
