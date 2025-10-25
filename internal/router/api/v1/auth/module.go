package auth

import (
	gonethttp "github.com/ralvarezdev/go-net/http"
	pbauth "github.com/ralvarezdev/grpc-auth-proto-go/compiled/ralvarezdev/auth"

	internalmiddleware "github.com/ralvarezdev/uru-mobiles-recipes-api/internal/middleware"
)

var (
	Controller = controller{}
	Module     = &gonethttp.Module{
		Pattern: "/auth",
		AddHandlersFn: func(m *gonethttp.Module) {
			m.AddEndpointHandler(
				"POST /signup",
				Controller.SignUp,
				internalmiddleware.Authenticate(
					pbauth.Auth_SignUp_FullMethodName,
				),
				internalmiddleware.ValidateProtoJSON(
					pbauth.SignUpRequest{},
				),
			)
			m.AddEndpointHandler(
				"POST /login",
				Controller.LogIn,
				internalmiddleware.Authenticate(
					pbauth.Auth_LogIn_FullMethodName,
				),
				internalmiddleware.ValidateProtoJSON(
					pbauth.LogInRequest{},
				),
			)
			m.AddEndpointHandler(
				"POST /refresh-token",
				Controller.RefreshToken,
				internalmiddleware.Authenticate(
					pbauth.Auth_RefreshToken_FullMethodName,
				),
			)
			m.AddEndpointHandler(
				"POST /logout",
				Controller.LogOut,
				internalmiddleware.Authenticate(
					pbauth.Auth_LogOut_FullMethodName,
				),
			)
			m.AddEndpointHandler(
				"GET /refresh-token",
				Controller.GetRefreshToken,
				internalmiddleware.Authenticate(
					pbauth.Auth_GetRefreshToken_FullMethodName,
				),
				internalmiddleware.ValidateProtoJSON(
					pbauth.GetRefreshTokenRequest{},
				),
			)
			m.AddEndpointHandler(
				"GET /refresh-tokens",
				Controller.ListRefreshTokens,
				internalmiddleware.Authenticate(
					pbauth.Auth_ListRefreshTokens_FullMethodName,
				),
			)
			m.AddEndpointHandler(
				"DELETE /refresh-token",
				Controller.RevokeRefreshToken,
				internalmiddleware.Authenticate(
					pbauth.Auth_RevokeRefreshToken_FullMethodName,
				),
				internalmiddleware.ValidateProtoJSON(
					pbauth.RevokeRefreshTokenRequest{},
				),
			)
			m.AddEndpointHandler(
				"DELETE /refresh-tokens",
				Controller.RevokeRefreshTokens,
				internalmiddleware.Authenticate(
					pbauth.Auth_RevokeRefreshTokens_FullMethodName,
				),
			)
			m.AddEndpointHandler(
				"POST /2fa/totp/generate",
				Controller.Generate2FATOTPUrl,
				internalmiddleware.Authenticate(
					pbauth.Auth_Generate2FATOTPUrl_FullMethodName,
				),
			)
			m.AddEndpointHandler(
				"POST /2fa/totp/verify",
				Controller.Verify2FATOTP,
				internalmiddleware.Authenticate(
					pbauth.Auth_Verify2FATOTP_FullMethodName,
				),
				internalmiddleware.ValidateProtoJSON(
					pbauth.Verify2FATOTPRequest{},
				),
			)
			m.AddEndpointHandler(
				"DELETE /2fa/totp",
				Controller.Revoke2FATOTP,
				internalmiddleware.Authenticate(
					pbauth.Auth_Revoke2FATOTP_FullMethodName,
				),
			)
			m.AddEndpointHandler(
				"PUT /password",
				Controller.ChangePassword,
				internalmiddleware.Authenticate(
					pbauth.Auth_ChangePassword_FullMethodName,
				),
				internalmiddleware.ValidateProtoJSON(
					pbauth.ChangePasswordRequest{},
				),
			)
			m.AddEndpointHandler(
				"POST /password/forgot",
				Controller.ForgotPassword,
				internalmiddleware.Authenticate(
					pbauth.Auth_ForgotPassword_FullMethodName,
				),
				internalmiddleware.ValidateProtoJSON(
					pbauth.ForgotPasswordRequest{},
				),
			)
			m.AddEndpointHandler(
				"POST /password/reset",
				Controller.ResetPassword,
				internalmiddleware.Authenticate(
					pbauth.Auth_ResetPassword_FullMethodName,
				),
				internalmiddleware.ValidateProtoJSON(
					pbauth.ResetPasswordRequest{},
				),
			)
			m.AddEndpointHandler(
				"PUT /email",
				Controller.ChangeEmail,
				internalmiddleware.Authenticate(
					pbauth.Auth_ChangeEmail_FullMethodName,
				),
				internalmiddleware.ValidateProtoJSON(pbauth.ChangeEmailRequest{}),
			)
			m.AddEndpointHandler(
				"POST /email/send-verification",
				Controller.SendEmailVerificationToken,
				internalmiddleware.Authenticate(
					pbauth.Auth_SendEmailVerificationToken_FullMethodName,
				),
			)
			m.AddEndpointHandler(
				"POST /email/verify",
				Controller.VerifyEmail,
				internalmiddleware.Authenticate(
					pbauth.Auth_VerifyEmail_FullMethodName,
				),
				internalmiddleware.ValidateProtoJSON(pbauth.VerifyEmailRequest{}),
			)
			m.AddEndpointHandler(
				"PUT /phone-number",
				Controller.ChangePhoneNumber,
				internalmiddleware.Authenticate(
					pbauth.Auth_ChangePhoneNumber_FullMethodName,
				),
				internalmiddleware.ValidateProtoJSON(pbauth.ChangePhoneNumberRequest{}),
			)
			m.AddEndpointHandler(
				"POST /phone-number/send-verification",
				Controller.SendPhoneNumberVerificationCode,
				internalmiddleware.Authenticate(
					pbauth.Auth_SendPhoneNumberVerificationCode_FullMethodName,
				),
			)
			m.AddEndpointHandler(
				"POST /phone-number/verify",
				Controller.VerifyPhoneNumber,
				internalmiddleware.Authenticate(
					pbauth.Auth_VerifyPhoneNumber_FullMethodName,
				),
				internalmiddleware.ValidateProtoJSON(pbauth.VerifyPhoneNumberRequest{}),
			)
			m.AddEndpointHandler(
				"POST /2fa/enable",
				Controller.EnableUser2FA,
				internalmiddleware.Authenticate(
					pbauth.Auth_EnableUser2FA_FullMethodName,
				),
				internalmiddleware.ValidateProtoJSON(pbauth.EnableUser2FARequest{}),
			)
			m.AddEndpointHandler(
				"POST /2fa/disable",
				Controller.DisableUser2FA,
				internalmiddleware.Authenticate(
					pbauth.Auth_DisableUser2FA_FullMethodName,
				),
				internalmiddleware.ValidateProtoJSON(pbauth.DisableUser2FARequest{}),
			)
			m.AddEndpointHandler(
				"POST /2fa/recovery-codes/regenerate",
				Controller.RegenerateUser2FARecoveryCodes,
				internalmiddleware.Authenticate(
					pbauth.Auth_RegenerateUser2FARecoveryCodes_FullMethodName,
				),
				internalmiddleware.ValidateProtoJSON(pbauth.RegenerateUser2FARecoveryCodesRequest{}),
			)
			m.AddEndpointHandler(
				"POST /2fa/email/send-code",
				Controller.SendUser2FAEmailCode,
				internalmiddleware.Authenticate(
					pbauth.Auth_SendUser2FAEmailCode_FullMethodName,
				),
				internalmiddleware.ValidateProtoJSON(pbauth.SendUser2FAEmailCodeRequest{}),
			)
		},
	}
)
