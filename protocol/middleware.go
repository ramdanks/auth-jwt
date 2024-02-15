package protocol

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const (
	MIDDLEWARE_CLAIMS_KEY string = "x-claims"
)

type Middleware struct {
	AccessTokenSecret []byte
}

// Middleware function to ensures the validity of authorization signature
// before forwards call to handler.
//
// Disclaimer, this is not ensuring the user is authorized. Even if the signatures is valid,
// the jwt might've been expired. Thus, the handler may receive an invalid/expired token.
func (m *Middleware) EnsureSignature(ctx *gin.Context) {
	m.verifyAndSetClaims(ctx, true)
}

// Middleware function to ensures the user is authorized to his claims.
// The jwt signature is verified and not expired.
func (m *Middleware) EnsureAuthorized(ctx *gin.Context) {
	m.verifyAndSetClaims(ctx, false)
}

func (m *Middleware) verifyAndSetClaims(
	ctx *gin.Context,
	allowExpired bool,
) {
	var err error
	var header Header

	err = ctx.BindHeader(&header)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	bearerToken := strings.Replace(header.Authorization, "Bearer ", "", -1)

	token, err := jwt.ParseWithClaims(
		bearerToken,
		&Claims{},
		func(t *jwt.Token) (interface{}, error) {
			return m.AccessTokenSecret, nil
		},
	)

	// guard checking
	switch allowExpired {
	case true:
		if err == nil && token.Valid {
			break
		}

		// here we need to make an exception.
		// if the jwt is invalid because it's expired, it will
		// produce an error code => jwt.ValidationErrorExpired.

		errV, ok := err.(*jwt.ValidationError)
		if !ok {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if errV.Errors != jwt.ValidationErrorExpired {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

	case false:
		if err != nil || !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.Set(MIDDLEWARE_CLAIMS_KEY, *claims)
	ctx.Next()
}
