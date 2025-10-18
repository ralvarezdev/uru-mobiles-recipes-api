package user

import (
	gonethttp "github.com/ralvarezdev/go-net/http"
	pbauth "github.com/ralvarezdev/grpc-auth-proto-go/compiled/ralvarezdev/auth"
	internalmiddleware "github.com/ralvarezdev/uru-mobiles-recipes-api/internal/middleware"
)

var (
	Controller = controller{}
	Module     = &gonethttp.Module{
		Pattern: "/user",
		AddHandlersFn: func(m *gonethttp.Module) {
			m.AddEndpointHandler(
				"PUT /profile",
				Controller.UpdateProfile,
				internalmiddleware.Authenticate(
					pbauth.Auth_UpdateProfile_FullMethodName,
				),
				internalmiddleware.ValidateProtoJSON(pbauth.UpdateProfileRequest{}),
			)
			m.AddEndpointHandler(
				"GET /profile",
				Controller.GetMyProfile,
				internalmiddleware.Authenticate(
					pbauth.Auth_GetMyProfile_FullMethodName,
				),
			)
			m.AddEndpointHandler(
				"PUT /username",
				Controller.ChangeUsername,
				internalmiddleware.Authenticate(
					pbauth.Auth_ChangeUsername_FullMethodName,
				),
				internalmiddleware.ValidateProtoJSON(pbauth.ChangeUsernameRequest{}),
			)
			m.AddEndpointHandler(
				"DELETE /",
				Controller.DeleteUser,
				internalmiddleware.Authenticate(
					pbauth.Auth_DeleteUser_FullMethodName,
				),
				internalmiddleware.ValidateProtoJSON(pbauth.DeleteUserRequest{}),
			)
		},
	}
)
