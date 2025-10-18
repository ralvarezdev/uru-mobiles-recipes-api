package user

import (
	"context"
	"net/http"

	gogrpcnethttp "github.com/ralvarezdev/go-grpc/client/net/http"
	gonethttpctx "github.com/ralvarezdev/go-net/http/context"
	gonethttpresponsejsend "github.com/ralvarezdev/go-net/http/response/jsend"
	gonethttpresponsejsendgrpc "github.com/ralvarezdev/go-net/http/response/jsend/grpc"
	pbauth "github.com/ralvarezdev/grpc-auth-proto-go/compiled/ralvarezdev/auth"
	internalgrpcauth "github.com/ralvarezdev/rest-auth/internal/grpc/auth"
	internaljson "github.com/ralvarezdev/rest-auth/internal/json"
	internalprotojson "github.com/ralvarezdev/rest-auth/internal/protojson"
	pbempty "google.golang.org/protobuf/types/known/emptypb"
)

type (
	// controller is the structure for the API V1 user controller
	controller struct{}
)

// UpdateProfile updates the profile of the authenticated user
// @Summary Updates the profile of the authenticated user
// @Description Updates the profile of the authenticated user
// @Tags api v1 user
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param request body pbauth.UpdateProfileRequest true "Update Profile Request"
// @Success 200 {object} gonethttpresponsejsend.SuccessBody[nil]
// @Failure 400 {object} gonethttpresponsejsend.FailBody
// @Failure 401 {object} gonethttpresponsejsend.FailBody
// @Failure 500 {object} gonethttpresponsejsend.ErrorBody
// @Router /api/v1/user/profile [put]
func (c controller) UpdateProfile(
	w http.ResponseWriter,
	r *http.Request,
) error {
	// Get the body from the context
	requestBody, _ := gonethttpctx.GetBody(r).(*pbauth.UpdateProfileRequest)

	// Create the context for the gRPC call
	ctx, err := gogrpcnethttp.SetCtxMetadataAuthorizationToken(
		context.Background(),
		r,
	)
	if err != nil {
		return err
	}

	// Call the gRPC service to update the profile
	if _, err = internalgrpcauth.Client.UpdateProfile(
		ctx,
		requestBody,
	); err != nil {
		return gonethttpresponsejsendgrpc.ParseError(err, true)
	}

	// Handle the response
	internaljson.Handler.HandleResponse(
		w, gonethttpresponsejsend.NewSuccessResponse(nil, http.StatusOK),
	)
	return nil
}

// GetMyProfile gets the profile of the authenticated user
// @Summary Gets the profile of the authenticated user
// @Description Gets the profile of the authenticated user
// @Tags api v1 user
// @Accept json
// @Produce json
// @Security CookieAuth
// @Success 200 {object} gonethttpresponsejsend.SuccessBody[pbauth.GetMyProfileResponse]
// @Failure 401 {object} gonethttpresponsejsend.FailBody
// @Failure 500 {object} gonethttpresponsejsend.ErrorBody
// @Router /api/v1/user/profile [get]
func (c controller) GetMyProfile(
	w http.ResponseWriter,
	r *http.Request,
) error {
	// Create the context for the gRPC call
	ctx, err := gogrpcnethttp.SetCtxMetadataAuthorizationToken(
		context.Background(),
		r,
	)
	if err != nil {
		return err
	}

	// Call the gRPC service to update the profile
	responseBody, err := internalgrpcauth.Client.GetMyProfile(
		ctx,
		&pbempty.Empty{},
	)
	if err != nil {
		return gonethttpresponsejsendgrpc.ParseError(err, true)
	}

	// Handle the response
	internalprotojson.Handler.HandleResponse(
		w,
		gonethttpresponsejsend.NewSuccessResponse(responseBody, http.StatusOK),
	)
	return nil
}

// ChangeUsername changes the username of the authenticated user
// @Summary Changes the username of the authenticated user
// @Description Changes the username of the authenticated user
// @Tags api v1 user
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param request body pbauth.ChangeUsernameRequest true "Change Username Request"
// @Success 200 {object} gonethttpresponsejsend.SuccessBody[nil]
// @Failure 400 {object} gonethttpresponsejsend.FailBody
// @Failure 401 {object} gonethttpresponsejsend.FailBody
// @Failure 500 {object} gonethttpresponsejsend.ErrorBody
// @Router /api/v1/user/username [put]
func (c controller) ChangeUsername(
	w http.ResponseWriter,
	r *http.Request,
) error {
	// Get the body from the context
	requestBody, _ := gonethttpctx.GetBody(r).(*pbauth.ChangeUsernameRequest)

	// Create the context for the gRPC call
	ctx, err := gogrpcnethttp.SetCtxMetadataAuthorizationToken(
		context.Background(),
		r,
	)
	if err != nil {
		return err
	}

	// Call the gRPC service to change the username
	if _, err = internalgrpcauth.Client.ChangeUsername(
		ctx,
		requestBody,
	); err != nil {
		return gonethttpresponsejsendgrpc.ParseError(err, true)
	}

	// Handle the response
	internaljson.Handler.HandleResponse(
		w,
		gonethttpresponsejsend.NewSuccessResponse(nil, http.StatusOK),
	)
	return nil
}

// DeleteUser deletes the authenticated user
// @Summary Deletes the authenticated user
// @Description Deletes the authenticated user
// @Tags api v1 user
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param request body pbauth.DeleteUserRequest true "Delete User Request"
// @Success 200 {object} gonethttpresponsejsend.SuccessBody[nil]
// @Failure 400 {object} gonethttpresponsejsend.FailBody
// @Failure 401 {object} gonethttpresponsejsend.FailBody
// @Failure 500 {object} gonethttpresponsejsend.ErrorBody
// @Router /api/v1/user [delete]
func (c controller) DeleteUser(
	w http.ResponseWriter,
	r *http.Request,
) error {
	// Get the body from the context
	requestBody, _ := gonethttpctx.GetBody(r).(*pbauth.DeleteUserRequest)

	// Create the context for the gRPC call
	ctx, err := gogrpcnethttp.SetCtxMetadataAuthorizationToken(
		context.Background(),
		r,
	)
	if err != nil {
		return err
	}

	// Call the gRPC service to delete the user
	if _, err = internalgrpcauth.Client.DeleteUser(
		ctx,
		requestBody,
	); err != nil {
		return gonethttpresponsejsendgrpc.ParseError(err, true)
	}

	// Handle the response
	internaljson.Handler.HandleResponse(
		w, gonethttpresponsejsend.NewSuccessResponse(nil, http.StatusOK),
	)
	return nil
}
