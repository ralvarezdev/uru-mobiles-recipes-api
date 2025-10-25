package router

import (
	gonethttp "github.com/ralvarezdev/go-net/http"
	gosecurityheadersnethttp "github.com/ralvarezdev/go-security-headers/net/http"

	internalmiddleware "github.com/ralvarezdev/uru-mobiles-recipes-api/internal/middleware"
	internalrouterapi "github.com/ralvarezdev/uru-mobiles-recipes-api/internal/router/api"
	internalrouterdocs "github.com/ralvarezdev/uru-mobiles-recipes-api/internal/router/docs"
)

var (
	Module = &gonethttp.Module{
		Pattern: "/",
		BeforeLoadFn: func(m *gonethttp.Module) {
			m.Middlewares = gonethttp.NewMiddlewares(
				internalmiddleware.LimitRequests,
				internalmiddleware.LimitBody,
				internalmiddleware.HandleError,
				gosecurityheadersnethttp.Handler,
			)
		},
		Submodules: gonethttp.NewSubmodules(
			internalrouterapi.Module,
			internalrouterdocs.Module,
		),
	}
)
