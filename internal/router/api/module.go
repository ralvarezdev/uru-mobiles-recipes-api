package api

import (
	gonethttp "github.com/ralvarezdev/go-net/http"
	internalrouterapiv1 "github.com/ralvarezdev/uru-mobiles-recipes-api/internal/router/api/v1"
)

var (
	Module = &gonethttp.Module{
		Pattern:    "/api",
		Submodules: gonethttp.NewSubmodules(internalrouterapiv1.Module),
	}
)
