package v1

import (
	"net/http"

	gonethttpresponsejsend "github.com/ralvarezdev/go-net/http/response/jsend"

	internaljson "github.com/ralvarezdev/uru-mobiles-recipes-api/internal/json"
)

// Ping pings the service
// @Summary Ping the service
// @Description Returns a pong response to check if the service is running
// @Tags api v1
// @Accept json
// @Produce json
// @Success 200 {object} gonethttpresponsejsend.SuccessBody[any]
// @Router /api/v1/ping [get]
func Ping(w http.ResponseWriter, r *http.Request) error {
	// Handle the response
	internaljson.Handler.HandleResponse(
		w, r, gonethttpresponsejsend.NewSuccessResponse(
			nil, http.StatusOK,
		),
	)
	return nil
}
