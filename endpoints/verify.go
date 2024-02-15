package endpoints

import (
	"example/auth-jwt/protocol"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewVerify() *protocol.Endpoint {
	return &protocol.Endpoint{
		Method:       protocol.GET,
		RelativePath: "/verify",
		Setup:        verifySetup,
	}
}

func verifySetup(
	env *protocol.Environment,
	mid *protocol.Middleware,
) ([]gin.HandlerFunc, protocol.ErrorCode) {
	return []gin.HandlerFunc{
		mid.EnsureAuthorized,
		verify,
	}, protocol.ERR_NONE
}

func verify(ctx *gin.Context) {
	claims := ctx.MustGet(protocol.MIDDLEWARE_CLAIMS_KEY).(protocol.Claims)
	message := fmt.Sprintf("%s, you are Authorized", claims.Username)
	ctx.String(http.StatusOK, message)
}
