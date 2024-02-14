package endpoints

import (
	"encoding/hex"
	"example/auth-jwt/protocol"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// All variables is a reference
// and will be loaded on server initialize
var (
	signingMethod         jwt.SigningMethod
	accessTokenSecret     []byte
	refreshTokenSecret    []byte
	accessExpirationTime  time.Duration
	refreshExpirationTime time.Duration
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

func NewLogin() *protocol.Endpoint {
	return &protocol.Endpoint{
		Method:       protocol.POST,
		RelativePath: "/login",
		Setup:        loginSetup,
	}
}

func loginSetup(
	env *protocol.Environment,
	mid *protocol.Middleware,
) ([]gin.HandlerFunc, protocol.ErrorCode) {
	var envval string
	var err error
	var ok bool

	accessTokenSecret = mid.AccessTokenSecret
	signingMethod = jwt.SigningMethodHS256

	envval, ok = os.LookupEnv(protocol.ENV_KEY_REFRESH_TOKEN)
	if !ok {
		return nil, protocol.ERR_ENV_REFRESH_TOKEN_NOT_SPECIFIED
	}

	refreshTokenSecret, err = hex.DecodeString(envval)
	if err != nil {
		return nil, protocol.ERR_ENV_REFRESH_TOKEN_MALFORMED
	}

	envval, ok = env.Lookup(protocol.ENV_KEY_ACCESS_EXP_TIME)
	if !ok {
		return nil, protocol.ERR_ENV_ACCESS_EXP_TIME_NOT_SPECIFIED
	}

	accessExpirationTime, err = time.ParseDuration(envval)
	if err != nil {
		return nil, protocol.ERR_ENV_ACCESS_EXP_TIME_MALFORMED
	}

	envval, ok = env.Lookup(protocol.ENV_KEY_REFRESH_EXP_TIME)
	if !ok {
		return nil, protocol.ERR_ENV_REFRESH_EXP_TIME_NOT_SPECIFIED
	}

	refreshExpirationTime, err = time.ParseDuration(envval)
	if err != nil {
		return nil, protocol.ERR_ENV_REFRESH_EXP_TIME_MALFORMED
	}

	return []gin.HandlerFunc{login}, protocol.ERR_NONE
}

func login(c *gin.Context) {
	var request loginRequest

	err := c.BindJSON(&request)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if invalidUsername(request.Username) || invalidPassword(request.Password) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ok := validCredentials(request)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	aToken := jwt.NewWithClaims(
		signingMethod,
		protocol.NewClaims(accessExpirationTime, request.Username),
	)

	rToken := jwt.NewWithClaims(
		signingMethod,
		protocol.NewClaims(refreshExpirationTime, request.Username),
	)

	aSigned, err := aToken.SignedString(accessTokenSecret)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	rSigned, err := rToken.SignedString(refreshTokenSecret)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, loginResponse{
		Access:  aSigned,
		Refresh: rSigned,
	})
}

var users = [...]loginRequest{
	{
		Username: "admin",
		Password: "admin-auth",
	},
	{
		Username: "ramadhanks",
		Password: "12345678",
	},
}

func invalidUsername(s string) bool {
	if len(s) < 4 {
		return true
	}

	for _, c := range s {
		if c == ' ' {
			return true
		}
	}
	return false
}

func invalidPassword(s string) bool {
	if len(s) < 8 {
		return true
	}

	for _, c := range s {
		if c == ' ' {
			return true
		}
	}
	return false
}

func validCredentials(p loginRequest) bool {
	for _, u := range users {
		if u.Username == p.Username {
			return u.Password == p.Password
		}
	}
	return false
}
