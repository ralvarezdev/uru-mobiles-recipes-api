package v1

import (
	gonethttp "github.com/ralvarezdev/go-net/http"

	internalrouterapiv1auth "github.com/ralvarezdev/uru-mobiles-recipes-api/internal/router/api/v1/auth"
	internalrouterapiv1user "github.com/ralvarezdev/uru-mobiles-recipes-api/internal/router/api/v1/user"
)

var (
	Module     = &gonethttp.Module{
		Pattern: "/v1",
		Submodules: gonethttp.NewSubmodules(
			internalrouterapiv1auth.Module,
			internalrouterapiv1user.Module,
		),
		AddHandlersFn: func(m *gonethttp.Module) {
			m.AddEndpointHandler(
				"GET /ping",
				Ping,
			)
		},
	}
)
