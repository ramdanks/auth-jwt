package endpoints

import (
	"example/auth-jwt/protocol"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewPing() *protocol.Endpoint {
	return &protocol.Endpoint{
		Method:       protocol.GET,
		RelativePath: "/ping",
		Setup:        pingSetup,
	}
}

func pingSetup(
	env *protocol.Environment,
	mid *protocol.Middleware,
) ([]gin.HandlerFunc, protocol.ErrorCode) {
	return []gin.HandlerFunc{ping}, protocol.ERR_NONE
}

func ping(c *gin.Context) {
	c.String(http.StatusOK, "Pong!")
}
