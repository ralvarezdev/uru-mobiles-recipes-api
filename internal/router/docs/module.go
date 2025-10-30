package router

import (
	"net/http"
	
	gonethttp "github.com/ralvarezdev/go-net/http"
	internaljson "github.com/ralvarezdev/uru-mobiles-recipes-api/internal/json"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

var (
	Module = &gonethttp.Module{
		Pattern: "/docs",
		AddHandlersFn: func(m *gonethttp.Module) {
			m.AddHandleFunc(
				"GET /swagger/",
				httpSwagger.Handler(
					httpSwagger.URL("./swagger.json"), // The URL pointing to API definition
					httpSwagger.DeepLinking(true),
				),
			)
			m.AddHandleFunc(
				"GET /swagger/swagger.json",
				func(w http.ResponseWriter, r *http.Request) {
    				w.Header().Set("Content-Type", "application/json")
      				w.Write(internaljson.SwaggerJSONDefinitions)
				},
			)
		},
	}
)
