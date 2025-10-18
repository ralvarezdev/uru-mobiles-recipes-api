package v1

import (
	"net/http"

	gonethttpresponsejsend "github.com/ralvarezdev/go-net/http/response/jsend"
	internaljson "github.com/ralvarezdev/rest-auth/internal/json"
)

type (
	// controller is the structure for the API V1 controller
	controller struct{}
)

// Ping pings the service
// @Summary Ping the service
// @Description Returns a pong response to check if the service is running
// @Tags api v1
// @Accept json
// @Produce json
// @Success 200 {object} gonethttpresponse.JSendSuccessBody
// @Router /api/v1/ping [get]
func (c controller) Ping(w http.ResponseWriter, r *http.Request) error {
	// Handle the response
	internaljson.Handler.HandleResponse(
		w, gonethttpresponsejsend.NewSuccessResponse(
			nil, http.StatusOK,
		),
	)
	return nil
}
