package auth

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
	// controller is the structure for the API V1 auth controller
	controller struct{}
)

// SignUp signs up a new user
// @Summary Sign up a new user
// @Description Creates a new user account with the provided details
// @Tags api v1 auth
// @Accept json
// @Produce json
// @Param request body pbauth.SignUpRequest true "Sign Up Request"
// @Success 201 {object} gonethttpresponsejsend.SuccessBody[pbauth.SignUpResponse]
// @Failure 400 {object} gonethttpresponsejsend.FailBody
// @Failure 500 {object} gonethttpresponsejsend.ErrorBody
// @Router /api/v1/auth/signup [post]
func (c controller) SignUp(w http.ResponseWriter, r *http.Request) error {
	// Get the body from the context
	requestBody, _ := gonethttpctx.GetBody(r).(*pbauth.SignUpRequest)

	// Call the gRPC service to sign up the user
	if _, err := internalgrpcauth.Client.SignUp(
		context.Background(),
		requestBody,
	); err != nil {
		return gonethttpresponsejsendgrpc.ParseError(err, true)
	}

	// Handle the response
	internaljson.Handler.HandleResponse(
		w, gonethttpresponsejsend.NewSuccessResponse(
			nil, http.StatusCreated,
		),
	)
	return nil
}

// LogIn logs in a user
// @Summary Log in a user
// @Description Authenticates a user and returns a seed token
// @Tags api v1 auth
// @Accept json
// @Produce json
// @Param request body pbauth.LogInRequest true "Log In Request"
// @Success 201 {object} LogInResponseBody[pbauth.LogInResponse]
// @Failure 400 {object} gonethttpresponsejsend.FailBody
// @Failure 401 {object} gonethttpresponsejsend.FailBody
// @Failure 500 {object} gonethttpresponsejsend.ErrorBody
// @Router /api/v1/auth/login [post]
func (c controller) LogIn(w http.ResponseWriter, r *http.Request) error {
	// Get the body from the context
	requestBody, _ := gonethttpctx.GetBody(r).(*pbauth.LogInRequest)

	// Call the gRPC service to log in the user
	ctx := context.Background()
	responseBody, err := internalgrpcauth.Client.LogIn(
		ctx,
		requestBody,
	)
	if err != nil {
		return gonethttpresponsejsendgrpc.ParseError(err, true)
	}

	// Handle the response (if the response contains 2FA methods, return a fail response)
	if responseBody != nil && responseBody.GetTwoFactorMethods() != nil {
		internalprotojson.Handler.HandleResponse(
			w, gonethttpresponsejsend.NewFailResponse(
				responseBody,
				http.StatusBadRequest,
			),
		)
		return nil
	}

	// Parse the metadata to clear cookies
	if err = internalgrpcauth.AuthenticationParser.ParseAuthorizationMetadataAsCookie(
		w,
		ctx,
	); err != nil {
		return err
	}

	// Handle the response
	internaljson.Handler.HandleResponse(
		w, gonethttpresponsejsend.NewSuccessResponse(
			nil,
			http.StatusCreated,
		),
	)
	return nil
}

// ListRefreshTokens gets a user's refresh tokens
// @Summary Get a user's refresh tokens
// @Description Gets a user's refresh tokens
// @Tags api v1 auth
// @Accept json
// @Produce json
// @Security CookieAuth
// @Success 200 {object} gonethttpresponsejsend.SuccessBody[pbauth.ListRefreshTokensResponse]
// @Failure 401 {object} gonethttpresponsejsend.FailBody
// @Failure 500 {object} gonethttpresponsejsend.ErrorBody
// @Router /api/v1/auth/refresh-tokens [get]
func (c controller) ListRefreshTokens(
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

	// Call the gRPC service to list the refresh tokens
	responseBody, err := internalgrpcauth.Client.ListRefreshTokens(
		ctx,
		&pbempty.Empty{},
	)
	if err != nil {
		return gonethttpresponsejsendgrpc.ParseError(err, true)
	}

	// Handle the response
	internalprotojson.Handler.HandleResponse(
		w, gonethttpresponsejsend.NewSuccessResponse(
			responseBody,
			http.StatusOK,
		),
	)
	return nil
}

// GetRefreshToken gets a user's refresh token
// @Summary Get a user's refresh token
// @Description Gets a user's refresh token
// @Tags api v1 auth
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param request body pbauth.GetRefreshTokenRequest true "Get Refresh Token Request"
// @Success 200 {object} gonethttpresponsejsend.SuccessBody[pbauth.GetRefreshTokenResponse]
// @Failure 401 {object} gonethttpresponsejsend.FailBody
// @Failure 500 {object} gonethttpresponsejsend.ErrorBody
// @Router /api/v1/auth/refresh-token [get]
func (c controller) GetRefreshToken(
	w http.ResponseWriter,
	r *http.Request,
) error {
	// Get the body from the context
	requestBody, _ := gonethttpctx.GetBody(r).(*pbauth.GetRefreshTokenRequest)

	// Create the context for the gRPC call
	ctx, err := gogrpcnethttp.SetCtxMetadataAuthorizationToken(
		context.Background(),
		r,
	)
	if err != nil {
		return err
	}

	// Call the gRPC service to get the refresh token
	responseBody, err := internalgrpcauth.Client.GetRefreshToken(
		ctx,
		requestBody,
	)
	if err != nil {
		return gonethttpresponsejsendgrpc.ParseError(err, true)
	}

	// Handle the response
	internalprotojson.Handler.HandleResponse(
		w, gonethttpresponsejsend.NewSuccessResponse(
			responseBody,
			http.StatusOK,
		),
	)
	return nil
}

// RevokeRefreshToken revokes a user's refresh token
// @Summary Revoke a user's refresh token
// @Description Revokes a user's refresh token
// @Tags api v1 auth
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param request body pbauth.RevokeRefreshTokenRequest true "Revoke Refresh Token Request"
// @Success 200 {object} gonethttpresponsejsend.SuccessBody[nil]
// @Failure 401 {object} gonethttpresponsejsend.FailBody
// @Failure 404 {object} gonethttpresponsejsend.FailBody
// @Failure 500 {object} gonethttpresponsejsend.ErrorBody
// @Router /api/v1/auth/refresh-token [delete]
func (c controller) RevokeRefreshToken(
	w http.ResponseWriter,
	r *http.Request,
) error {
	// Get the body from the context
	requestBody, _ := gonethttpctx.GetBody(r).(*pbauth.RevokeRefreshTokenRequest)

	// Create the context for the gRPC call
	ctx, err := gogrpcnethttp.SetCtxMetadataAuthorizationToken(
		context.Background(),
		r,
	)
	if err != nil {
		return err
	}

	// Call the gRPC service to revoke the refresh token
	if _, err = internalgrpcauth.Client.RevokeRefreshToken(
		ctx,
		requestBody,
	); err != nil {
		return gonethttpresponsejsendgrpc.ParseError(err, true)
	}

	// Parse the metadata to clear cookies
	if err = internalgrpcauth.AuthenticationParser.ParseAuthorizationMetadataAsCookie(
		w,
		ctx,
	); err != nil {
		return err
	}

	// Handle the response
	internaljson.Handler.HandleResponse(
		w, gonethttpresponsejsend.NewSuccessResponse(
			nil,
			http.StatusOK,
		),
	)
	return nil
}

// LogOut logs out a user
// @Summary Log out a user
// @Description Logs out a user
// @Tags api v1 auth
// @Accept json
// @Produce json
// @Security CookieAuth
// @Success 200 {object} gonethttpresponsejsend.SuccessBody[nil]
// @Failure 401 {object} gonethttpresponsejsend.FailBody
// @Failure 500 {object} gonethttpresponsejsend.ErrorBody
// @Router /api/v1/auth/logout [post]
func (c controller) LogOut(w http.ResponseWriter, r *http.Request) error {
	// Create the context for the gRPC call
	ctx, err := gogrpcnethttp.SetCtxMetadataAuthorizationToken(
		context.Background(),
		r,
	)
	if err != nil {
		return err
	}

	// Call the gRPC service to log out
	if _, err = internalgrpcauth.Client.LogOut(
		ctx,
		&pbempty.Empty{},
	); err != nil {
		return gonethttpresponsejsendgrpc.ParseError(err, true)
	}

	// Parse the metadata to clear cookies
	if err = internalgrpcauth.AuthenticationParser.ParseAuthorizationMetadataAsCookie(
		w,
		ctx,
	); err != nil {
		return err
	}

	// Handle the response
	internaljson.Handler.HandleResponse(
		w, gonethttpresponsejsend.NewSuccessResponse(
			nil,
			http.StatusOK,
		),
	)
	return nil
}

// RevokeRefreshTokens revokes a user's refresh tokens
// @Summary Revoke a user's refresh tokens
// @Description Revokes a user's refresh tokens
// @Tags api v1 auth
// @Accept json
// @Produce json
// @Security CookieAuth
// @Success 200 {object} gonethttpresponsejsend.SuccessBody[nil]
// @Failure 401 {object} gonethttpresponsejsend.FailBody
// @Failure 500 {object} gonethttpresponsejsend.ErrorBody
// @Router /api/v1/auth/refresh-tokens [delete]
func (c controller) RevokeRefreshTokens(
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

	// Call the gRPC service to revoke the refresh tokens
	if _, err = internalgrpcauth.Client.RevokeRefreshTokens(
		ctx,
		&pbempty.Empty{},
	); err != nil {
		return gonethttpresponsejsendgrpc.ParseError(err, true)
	}

	// Parse the metadata to clear cookies
	if err = internalgrpcauth.AuthenticationParser.ParseAuthorizationMetadataAsCookie(
		w,
		ctx,
	); err != nil {
		return err
	}

	// Handle the response
	internaljson.Handler.HandleResponse(
		w, gonethttpresponsejsend.NewSuccessResponse(
			nil,
			http.StatusOK,
		),
	)
	return nil
}

// RefreshToken refreshes a user token
// @Summary Refresh a user token
// @Description Refreshes a user token
// @Tags api v1 auth
// @Accept json
// @Produce json
// @Security CookieAuth
// @Success 201 {object} gonethttpresponsejsend.SuccessBody[nil]
// @Failure 401 {object} gonethttpresponsejsend.FailBody
// @Failure 500 {object} gonethttpresponsejsend.ErrorBody
// @Router /api/v1/auth/refresh-token [post]
func (c controller) RefreshToken(w http.ResponseWriter, r *http.Request) error {
	// Create the context for the gRPC call
	ctx, err := gogrpcnethttp.SetCtxMetadataAuthorizationToken(
		context.Background(),
		r,
	)
	if err != nil {
		return err
	}

	// Call the gRPC service to refresh the token
	if _, err = internalgrpcauth.Client.RefreshToken(
		ctx,
		&pbempty.Empty{},
	); err != nil {
		return gonethttpresponsejsendgrpc.ParseError(err, true)
	}

	// Parse the metadata to clear cookies
	if err = internalgrpcauth.AuthenticationParser.ParseAuthorizationMetadataAsCookie(
		w,
		ctx,
	); err != nil {
		return err
	}

	// Handle the response
	internaljson.Handler.HandleResponse(
		w, gonethttpresponsejsend.NewSuccessResponse(
			nil,
			http.StatusCreated,
		),
	)
	return nil
}

// Generate2FATOTPUrl generates a 2FA TOTP URL
// @Summary Generate a 2FA TOTP URL
// @Description Generates a 2FA TOTP URL
// @Tags api v1 auth
// @Accept json
// @Produce json
// @Security CookieAuth
// @Success 201 {object} gonethttpresponsejsend.SuccessBody[pbauth.Generate2FATOTPUrlResponse]
// @Failure 401 {object} gonethttpresponsejsend.FailBody
// @Failure 500 {object} gonethttpresponsejsend.ErrorBody
// @Router /api/v1/auth/2fa/totp/generate [post]
func (c controller) Generate2FATOTPUrl(
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

	// Call the gRPC service to generate the 2FA TOTP URL
	responseBody, err := internalgrpcauth.Client.Generate2FATOTPUrl(
		ctx,
		&pbempty.Empty{},
	)
	if err != nil {
		return gonethttpresponsejsendgrpc.ParseError(err, true)
	}

	// Handle the response
	internalprotojson.Handler.HandleResponse(
		w, gonethttpresponsejsend.NewSuccessResponse(
			responseBody,
			http.StatusCreated,
		),
	)
	return nil
}

// Verify2FATOTP verifies a 2FA TOTP code
// @Summary Verify a 2FA TOTP code
// @Description Verifies a 2FA TOTP code
// @Tags api v1 auth
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param request body pbauth.Verify2FATOTPRequest true "Verify 2FA TOTP Request"
// @Success 200 {object} gonethttpresponsejsend.SuccessBody[nil]
// @Failure 400 {object} gonethttpresponsejsend.FailBody
// @Failure 401 {object} gonethttpresponsejsend.FailBody
// @Failure 500 {object} gonethttpresponsejsend.ErrorBody
// @Router /api/v1/auth/2fa/totp/verify [post]
func (c controller) Verify2FATOTP(
	w http.ResponseWriter,
	r *http.Request,
) error {
	// Get the body from the context
	requestBody, _ := gonethttpctx.GetBody(r).(*pbauth.Verify2FATOTPRequest)

	// Create the context for the gRPC call
	ctx, err := gogrpcnethttp.SetCtxMetadataAuthorizationToken(
		context.Background(),
		r,
	)
	if err != nil {
		return err
	}

	// Call the gRPC service to verify the 2FA TOTP code
	if _, err = internalgrpcauth.Client.Verify2FATOTP(
		ctx,
		requestBody,
	); err != nil {
		return gonethttpresponsejsendgrpc.ParseError(err, true)
	}

	// Handle the response
	internaljson.Handler.HandleResponse(
		w, gonethttpresponsejsend.NewSuccessResponse(
			nil,
			http.StatusOK,
		),
	)
	return nil
}

// Revoke2FATOTP revokes a user's 2FA TOTP
// @Summary Revoke a user's 2FA TOTP
// @Description Revokes a user's 2FA TOTP
// @Tags api v1 auth
// @Accept json
// @Produce json
// @Security CookieAuth
// @Success 200 {object} gonethttpresponsejsend.SuccessBody[nil]
// @Failure 401 {object} gonethttpresponsejsend.FailBody
// @Failure 500 {object} gonethttpresponsejsend.ErrorBody
// @Router /api/v1/auth/2fa/totp [delete]
func (c controller) Revoke2FATOTP(
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

	// Call the gRPC service to revoke the 2FA TOTP
	if _, err = internalgrpcauth.Client.Revoke2FATOTP(
		ctx,
		&pbempty.Empty{},
	); err != nil {
		return gonethttpresponsejsendgrpc.ParseError(err, true)
	}

	// Handle the response
	internaljson.Handler.HandleResponse(
		w, gonethttpresponsejsend.NewSuccessResponse(
			nil,
			http.StatusOK,
		),
	)
	return nil
}

// ChangeEmail changes the email of the authenticated user
// @Summary Changes the email of the authenticated user
// @Description Changes the email of the authenticated user
// @Tags api v1 auth
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param request body pbauth.ChangeEmailRequest true "Change Email Request"
// @Success 200 {object} gonethttpresponsejsend.SuccessBody[nil]
// @Failure 400 {object} gonethttpresponsejsend.FailBody
// @Failure 401 {object} gonethttpresponsejsend.FailBody
// @Failure 500 {object} gonethttpresponsejsend.ErrorBody
// @Router /api/v1/auth/email [put]
func (c controller) ChangeEmail(
	w http.ResponseWriter,
	r *http.Request,
) error {
	// Get the body from the context
	requestBody, _ := gonethttpctx.GetBody(r).(*pbauth.ChangeEmailRequest)

	// Create the context for the gRPC call
	ctx, err := gogrpcnethttp.SetCtxMetadataAuthorizationToken(
		context.Background(),
		r,
	)
	if err != nil {
		return err
	}

	// Call the gRPC service to change the email
	if _, err = internalgrpcauth.Client.ChangeEmail(
		ctx,
		requestBody,
	); err != nil {
		return gonethttpresponsejsendgrpc.ParseError(err, true)
	}

	// Handle the response
	internaljson.Handler.HandleResponse(
		w, gonethttpresponsejsend.NewSuccessResponse(
			nil,
			http.StatusOK,
		),
	)
	return nil
}

// SendEmailVerificationToken sends an email verification token to the authenticated user
// @Summary Sends an email verification token to the authenticated user
// @Description Sends an email verification token to the authenticated user
// @Tags api v1 auth
// @Accept json
// @Produce json
// @Security CookieAuth
// @Success 200 {object} gonethttpresponsejsend.SuccessBody[nil]
// @Failure 400 {object} gonethttpresponsejsend.FailBody
// @Failure 401 {object} gonethttpresponsejsend.FailBody
// @Failure 500 {object} gonethttpresponsejsend.ErrorBody
// @Router /api/v1/auth/email/send-verification [post]
func (c controller) SendEmailVerificationToken(
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

	// Call the gRPC service to send the email verification token
	if _, err = internalgrpcauth.Client.SendEmailVerificationToken(
		ctx,
		&pbempty.Empty{},
	); err != nil {
		return gonethttpresponsejsendgrpc.ParseError(err, true)
	}

	// Handle the response
	internaljson.Handler.HandleResponse(
		w, gonethttpresponsejsend.NewSuccessResponse(
			nil,
			http.StatusOK,
		),
	)
	return nil
}

// VerifyEmail verifies the email of the authenticated user
// @Summary Verifies the email of the authenticated user
// @Description Verifies the email of the authenticated user
// @Tags api v1 auth
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param request body pbauth.VerifyEmailRequest true "Verify Email Request"
// @Success 200 {object} gonethttpresponsejsend.SuccessBody[nil]
// @Failure 400 {object} gonethttpresponsejsend.FailBody
// @Failure 500 {object} gonethttpresponsejsend.ErrorBody
// @Router /api/v1/auth/email/verify [post]
func (c controller) VerifyEmail(
	w http.ResponseWriter,
	r *http.Request,
) error {
	// Get the body from the context
	requestBody, _ := gonethttpctx.GetBody(r).(*pbauth.VerifyEmailRequest)

	// Create the context for the gRPC call
	ctx, err := gogrpcnethttp.SetCtxMetadataAuthorizationToken(
		context.Background(),
		r,
	)
	if err != nil {
		return err
	}

	// Call the gRPC service to verify the email
	if _, err = internalgrpcauth.Client.VerifyEmail(
		ctx,
		requestBody,
	); err != nil {
		return gonethttpresponsejsendgrpc.ParseError(err, true)
	}

	// Handle the response
	internaljson.Handler.HandleResponse(
		w, gonethttpresponsejsend.NewSuccessResponse(
			nil,
			http.StatusOK,
		),
	)
	return nil
}

// ChangePassword changes a user's password
// @Summary Change a user's password
// @Description Changes a user's password
// @Tags api v1 auth
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param request body pbauth.ChangePasswordRequest true "Change Password Request"
// @Success 200 {object} gonethttpresponsejsend.SuccessBody[nil]
// @Failure 400 {object} gonethttpresponsejsend.FailBody
// @Failure 401 {object} gonethttpresponsejsend.FailBody
// @Failure 500 {object} gonethttpresponsejsend.ErrorBody
// @Router /api/v1/auth/password [put]
func (c controller) ChangePassword(
	w http.ResponseWriter,
	r *http.Request,
) error {
	// Get the body from the context
	requestBody, _ := gonethttpctx.GetBody(r).(*pbauth.ChangePasswordRequest)

	// Create the context for the gRPC call
	ctx, err := gogrpcnethttp.SetCtxMetadataAuthorizationToken(
		context.Background(),
		r,
	)
	if err != nil {
		return err
	}

	// Call the gRPC service to change the password
	if _, err = internalgrpcauth.Client.ChangePassword(
		ctx,
		requestBody,
	); err != nil {
		return gonethttpresponsejsendgrpc.ParseError(err, true)
	}

	// Handle the response
	internaljson.Handler.HandleResponse(
		w, gonethttpresponsejsend.NewSuccessResponse(
			nil,
			http.StatusOK,
		),
	)
	return nil
}

// ForgotPassword sends a password reset email
// @Summary Send a password reset email
// @Description Sends a password reset email
// @Tags api v1 auth
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param request body pbauth.ForgotPasswordRequest true "Forgot Password Request"
// @Success 200 {object} gonethttpresponsejsend.SuccessBody[nil]
// @Failure 400 {object} gonethttpresponsejsend.FailBody
// @Failure 401 {object} gonethttpresponsejsend.FailBody
// @Failure 500 {object} gonethttpresponsejsend.ErrorBody
// @Router /api/v1/auth/password/forgot [post]
func (c controller) ForgotPassword(
	w http.ResponseWriter,
	r *http.Request,
) error {
	// Get the body from the context
	requestBody, _ := gonethttpctx.GetBody(r).(*pbauth.ForgotPasswordRequest)

	// Create the context for the gRPC call
	ctx, err := gogrpcnethttp.SetCtxMetadataAuthorizationToken(
		context.Background(),
		r,
	)
	if err != nil {
		return err
	}

	// Call the gRPC service to send the password reset email
	if _, err = internalgrpcauth.Client.ForgotPassword(
		ctx,
		requestBody,
	); err != nil {
		return gonethttpresponsejsendgrpc.ParseError(err, true)
	}

	// Handle the response
	internaljson.Handler.HandleResponse(
		w, gonethttpresponsejsend.NewSuccessResponse(
			nil,
			http.StatusOK,
		),
	)
	return nil
}

// ResetPassword resets a user's password
// @Summary Reset a user's password
// @Description Resets a user's password
// @Tags api v1 auth
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param request body pbauth.ResetPasswordRequest true "Reset Password Request"
// @Success 200 {object} gonethttpresponsejsend.SuccessBody[nil]
// @Failure 400 {object} gonethttpresponsejsend.FailBody
// @Failure 500 {object} gonethttpresponsejsend.ErrorBody
// @Router /api/v1/auth/password/reset [post]
func (c controller) ResetPassword(
	w http.ResponseWriter,
	r *http.Request,
) error {
	// Get the body from the context
	requestBody, _ := gonethttpctx.GetBody(r).(*pbauth.ResetPasswordRequest)

	// Create the context for the gRPC call
	ctx, err := gogrpcnethttp.SetCtxMetadataAuthorizationToken(
		context.Background(),
		r,
	)
	if err != nil {
		return err
	}

	// Call the gRPC service to reset the password
	if _, err = internalgrpcauth.Client.ResetPassword(
		ctx,
		requestBody,
	); err != nil {
		return gonethttpresponsejsendgrpc.ParseError(err, true)
	}

	// Handle the response
	internaljson.Handler.HandleResponse(
		w, gonethttpresponsejsend.NewSuccessResponse(
			nil,
			http.StatusOK,
		),
	)
	return nil
}

// ChangePhoneNumber changes the phone number of the authenticated user
// @Summary Changes the phone number of the authenticated user
// @Description Changes the phone number of the authenticated user
// @Tags api v1 auth
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param request body pbauth.ChangePhoneNumberRequest true "Change Phone Number Request"
// @Success 200 {object} gonethttpresponsejsend.SuccessBody[nil]
// @Failure 400 {object} gonethttpresponsejsend.FailBody
// @Failure 401 {object} gonethttpresponsejsend.FailBody
// @Failure 500 {object} gonethttpresponsejsend.ErrorBody
// @Router /api/v1/auth/phone-number [put]
func (c controller) ChangePhoneNumber(
	w http.ResponseWriter,
	r *http.Request,
) error {
	// Get the body from the context
	requestBody, _ := gonethttpctx.GetBody(r).(*pbauth.ChangePhoneNumberRequest)

	// Create the context for the gRPC call
	ctx, err := gogrpcnethttp.SetCtxMetadataAuthorizationToken(
		context.Background(),
		r,
	)
	if err != nil {
		return err
	}

	// Call the gRPC service to change the phone number
	if _, err = internalgrpcauth.Client.ChangePhoneNumber(
		ctx,
		requestBody,
	); err != nil {
		return gonethttpresponsejsendgrpc.ParseError(err, true)
	}

	// Handle the response
	internaljson.Handler.HandleResponse(
		w, gonethttpresponsejsend.NewSuccessResponse(
			nil,
			http.StatusOK,
		),
	)
	return nil
}

// SendPhoneNumberVerificationCode sends a phone number verification code to the authenticated user
// @Summary Sends a phone number verification code to the authenticated user
// @Description Sends a phone number verification code to the authenticated user
// @Tags api v1 auth
// @Accept json
// @Produce json
// @Security CookieAuth
// @Success 200 {object} gonethttpresponsejsend.SuccessBody[nil]
// @Failure 400 {object} gonethttpresponsejsend.FailBody
// @Failure 401 {object} gonethttpresponsejsend.FailBody
// @Failure 500 {object} gonethttpresponsejsend.ErrorBody
// @Router /api/v1/auth/phone-number/send-verification [post]
func (c controller) SendPhoneNumberVerificationCode(
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

	// Call the gRPC service to send the phone number verification code
	if _, err = internalgrpcauth.Client.SendPhoneNumberVerificationCode(
		ctx,
		&pbempty.Empty{},
	); err != nil {
		return gonethttpresponsejsendgrpc.ParseError(err, true)
	}

	// Handle the response
	internaljson.Handler.HandleResponse(
		w, gonethttpresponsejsend.NewSuccessResponse(
			nil,
			http.StatusOK,
		),
	)
	return nil
}

// VerifyPhoneNumber verifies the phone number of the authenticated user
// @Summary Verifies the phone number of the authenticated user
// @Description Verifies the phone number of the authenticated user
// @Tags api v1 auth
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param request body pbauth.VerifyPhoneNumberRequest true "Verify Phone Number Request"
// @Success 200 {object} gonethttpresponsejsend.SuccessBody[nil]
// @Failure 400 {object} gonethttpresponsejsend.FailBody
// @Failure 401 {object} gonethttpresponsejsend.FailBody
// @Failure 500 {object} gonethttpresponsejsend.ErrorBody
// @Router /api/v1/auth/phone-number/verify [post]
func (c controller) VerifyPhoneNumber(
	w http.ResponseWriter,
	r *http.Request,
) error {
	// Get the body from the context
	requestBody, _ := gonethttpctx.GetBody(r).(*pbauth.EnableUser2FARequest)

	// Create the context for the gRPC call
	ctx, err := gogrpcnethttp.SetCtxMetadataAuthorizationToken(
		context.Background(),
		r,
	)
	if err != nil {
		return err
	}

	// Call the gRPC service to enable 2FA for the user
	responseBody, err := internalgrpcauth.Client.EnableUser2FA(
		ctx,
		requestBody,
	)
	if err != nil {
		return gonethttpresponsejsendgrpc.ParseError(err, true)
	}

	// Handle the response
	internalprotojson.Handler.HandleResponse(
		w, gonethttpresponsejsend.NewSuccessResponse(
			responseBody,
			http.StatusOK,
		),
	)
	return nil
}

// EnableUser2FA enables 2FA for the authenticated user
// @Summary Enable 2FA for the authenticated user
// @Description Enables 2FA for the authenticated user
// @Tags api v1 auth
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param request body pbauth.EnableUser2FARequest true "Enable User 2FA Request"
// @Success 200 {object} gonethttpresponsejsend.SuccessBody[pbauth.EnableUser2FAResponse]
// @Failure 400 {object} gonethttpresponsejsend.FailBody
// @Failure 401 {object} gonethttpresponsejsend.FailBody
// @Failure 500 {object} gonethttpresponsejsend.ErrorBody
// @Router /api/v1/auth/2fa/enable [post]
func (c controller) EnableUser2FA(
	w http.ResponseWriter,
	r *http.Request,
) error {
	// Get the body from the context
	requestBody, _ := gonethttpctx.GetBody(r).(*pbauth.EnableUser2FARequest)

	// Create the context for the gRPC call
	ctx, err := gogrpcnethttp.SetCtxMetadataAuthorizationToken(
		context.Background(),
		r,
	)
	if err != nil {
		return err
	}

	// Call the gRPC service to enable 2FA for the user
	responseBody, err := internalgrpcauth.Client.EnableUser2FA(
		ctx,
		requestBody,
	)
	if err != nil {
		return gonethttpresponsejsendgrpc.ParseError(err, true)
	}

	// Handle the response
	internalprotojson.Handler.HandleResponse(
		w, gonethttpresponsejsend.NewSuccessResponse(
			responseBody,
			http.StatusOK,
		),
	)
	return nil
}

// DisableUser2FA disables 2FA for the authenticated user
// @Summary Disable 2FA for the authenticated user
// @Description Disables 2FA for the authenticated user
// @Tags api v1 auth
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param request body pbauth.DisableUser2FARequest true "Disable User 2FA Request"
// @Success 200 {object} gonethttpresponsejsend.SuccessBody[nil]
// @Failure 400 {object} gonethttpresponsejsend.FailBody
// @Failure 401 {object} gonethttpresponsejsend.FailBody
// @Failure 500 {object} gonethttpresponsejsend.ErrorBody
// @Router /api/v1/auth/2fa/disable [post]
func (c controller) DisableUser2FA(
	w http.ResponseWriter,
	r *http.Request,
) error {
	// Get the body from the context
	requestBody, _ := gonethttpctx.GetBody(r).(*pbauth.DisableUser2FARequest)

	// Create the context for the gRPC call
	ctx, err := gogrpcnethttp.SetCtxMetadataAuthorizationToken(
		context.Background(),
		r,
	)
	if err != nil {
		return err
	}

	// Call the gRPC service to disable 2FA for the user
	if _, err = internalgrpcauth.Client.DisableUser2FA(
		ctx,
		requestBody,
	); err != nil {
		return gonethttpresponsejsendgrpc.ParseError(err, true)
	}

	// Handle the response
	internaljson.Handler.HandleResponse(
		w, gonethttpresponsejsend.NewSuccessResponse(
			nil,
			http.StatusOK,
		),
	)
	return nil
}

// RegenerateUser2FARecoveryCodes regenerates the 2FA recovery codes for the authenticated user
// @Summary Regenerate 2FA recovery codes for the authenticated user
// @Description Regenerates the 2FA recovery codes for the authenticated user
// @Tags api v1 auth
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param request body pbauth.RegenerateUser2FARecoveryCodesRequest true "Regenerate User 2FA Recovery Codes Request"
// @Success 200 {object} gonethttpresponsejsend.SuccessBody[pbauth.RegenerateUser2FARecoveryCodesResponse]
// @Failure 400 {object} gonethttpresponsejsend.FailBody
// @Failure 401 {object} gonethttpresponsejsend.FailBody
// @Failure 500 {object} gonethttpresponsejsend.ErrorBody
// @Router /api/v1/auth/2fa/recovery-codes/regenerate [post]
func (c controller) RegenerateUser2FARecoveryCodes(
	w http.ResponseWriter,
	r *http.Request,
) error {
	// Get the body from the context
	requestBody, _ := gonethttpctx.GetBody(r).(*pbauth.RegenerateUser2FARecoveryCodesRequest)

	// Create the context for the gRPC call
	ctx, err := gogrpcnethttp.SetCtxMetadataAuthorizationToken(
		context.Background(),
		r,
	)
	if err != nil {
		return err
	}

	// Call the gRPC service to regenerate the 2FA recovery codes
	responseBody, err := internalgrpcauth.Client.RegenerateUser2FARecoveryCodes(
		ctx,
		requestBody,
	)
	if err != nil {
		return gonethttpresponsejsendgrpc.ParseError(err, true)
	}

	// Handle the response
	internalprotojson.Handler.HandleResponse(
		w, gonethttpresponsejsend.NewSuccessResponse(
			responseBody,
			http.StatusOK,
		),
	)
	return nil
}

// SendUser2FAEmailCode sends a 2FA email code
// @Summary Send 2FA email code
// @Description Sends a 2FA email code
// @Tags api v1 auth
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param request body pbauth.SendUser2FAEmailCodeRequest true "Send User 2FA Email Code Request"
// @Success 200 {object} gonethttpresponsejsend.SuccessBody[nil]
// @Failure 400 {object} gonethttpresponsejsend.FailBody
// @Failure 401 {object} gonethttpresponsejsend.FailBody
// @Failure 500 {object} gonethttpresponsejsend.ErrorBody
// @Router /api/v1/auth/2fa/email/send-code [post]
func (c controller) SendUser2FAEmailCode(
	w http.ResponseWriter,
	r *http.Request,
) error {
	// Get the body from the context
	requestBody, _ := gonethttpctx.GetBody(r).(*pbauth.SendUser2FAEmailCodeRequest)

	// Create the context for the gRPC call
	ctx, err := gogrpcnethttp.SetCtxMetadataAuthorizationToken(
		context.Background(),
		r,
	)
	if err != nil {
		return err
	}

	// Call the gRPC service to send the 2FA email code
	if _, err = internalgrpcauth.Client.SendUser2FAEmailCode(
		ctx,
		requestBody,
	); err != nil {
		return gonethttpresponsejsendgrpc.ParseError(err, true)
	}

	// Handle the response
	internaljson.Handler.HandleResponse(
		w, gonethttpresponsejsend.NewSuccessResponse(
			nil,
			http.StatusOK,
		),
	)
	return nil
}
