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

func (m *Middleware) VerifyAllowExpired(ctx *gin.Context) {
	m.verifyAndSetClaims(ctx, true)
}

func (m *Middleware) Verify(ctx *gin.Context) {
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
