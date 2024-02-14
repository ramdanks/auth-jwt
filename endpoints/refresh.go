package endpoints

import (
	"example/auth-jwt/protocol"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type refreshRequest struct {
	Refresh string `json:"refresh"`
}

type refreshResponse struct {
	Access string `json:"access"`
}

func NewRefresh() *protocol.Endpoint {
	return &protocol.Endpoint{
		Method:       protocol.POST,
		RelativePath: "/refresh",
		Setup:        refreshSetup,
	}
}

func refreshSetup(
	env *protocol.Environment,
	mid *protocol.Middleware,
) ([]gin.HandlerFunc, protocol.ErrorCode) {
	return []gin.HandlerFunc{
		mid.VerifyAllowExpired,
		refresh,
	}, protocol.ERR_NONE
}

func refresh(c *gin.Context) {
	var r refreshRequest

	err := c.BindJSON(&r)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	token, _ := jwt.ParseWithClaims(
		r.Refresh,
		&protocol.Claims{},
		func(t *jwt.Token) (interface{}, error) {
			return refreshTokenSecret, nil
		},
	)

	if !token.Valid {
		c.Status(http.StatusUnauthorized)
		return
	}

	rclaims, ok := token.Claims.(*protocol.Claims)
	if !ok {
		c.Status(http.StatusBadRequest)
		return
	}

	claims, ok := c.MustGet(protocol.MIDDLEWARE_CLAIMS_KEY).(protocol.Claims)
	if !ok {
		c.Status(http.StatusInternalServerError)
		return
	}

	// the username in authorization claims should be
	// the same with the claims in the refresh token
	if claims.Username != rclaims.Username {
		c.Status(http.StatusUnauthorized)
		return
	}

	// client can only issues a new access token if
	// the previous access token is expired
	if claims.Valid() == nil {
		c.Status(http.StatusTooEarly)
		return
	}

	// process to create a new access token with claims
	aToken := jwt.NewWithClaims(
		signingMethod,
		protocol.NewClaims(accessExpirationTime, claims.Username),
	)

	aSigned, err := aToken.SignedString(accessTokenSecret)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, refreshResponse{Access: aSigned})
}
