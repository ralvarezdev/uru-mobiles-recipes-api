package docs

import (
	"net/http"

	gonethttp "github.com/ralvarezdev/go-net/http"
	gonethttpresponsejsend "github.com/ralvarezdev/go-net/http/response/jsend"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	internaljson "github.com/ralvarezdev/uru-mobiles-recipes-api/internal/json"
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
					if _, err := w.Write(internaljson.SwaggerJSONDefinitions); err != nil {
						internaljson.Handler.HandleResponse(
							w, r,
							gonethttpresponsejsend.NewDebugErrorResponse(
								gonethttp.ErrInternalServerError.Error(),
								err.Error(),
								http.StatusInternalServerError,
							),
						)
					}
				},
			)
		},
	}
)
