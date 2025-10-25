package cookie

import (
	goflagsmode "github.com/ralvarezdev/go-flags/mode"
	gojwttoken "github.com/ralvarezdev/go-jwt/token"
	gonethttpcookie "github.com/ralvarezdev/go-net/http/cookie"
)

var (
	// AccessToken is the cookies attributes for the access token cookie
	AccessToken *gonethttpcookie.Attributes

	// RefreshToken is the cookies attributes for the refresh token cookie
	RefreshToken *gonethttpcookie.Attributes
)

// Load loads the cookie constants based on the mode
//
// Parameters:
//
//   - mode: the go-flags mode flag to determine if the environment is in debug mode
func Load(mode *goflagsmode.Flag) {
	// Determine if the cookies should be secure
	secure := mode.IsProd()

	// AccessToken is the cookies attributes for the access token cookie
	AccessToken = &gonethttpcookie.Attributes{
		Name:     gojwttoken.AccessToken.String(),
		HTTPOnly: true,
		Secure:   secure,
		Path:     "/",
	}

	// RefreshToken is the cookies attributes for the refresh token cookie
	RefreshToken = &gonethttpcookie.Attributes{
		Name:     gojwttoken.RefreshToken.String(),
		HTTPOnly: true,
		Secure:   secure,
		Path:     "/",
	}
}
