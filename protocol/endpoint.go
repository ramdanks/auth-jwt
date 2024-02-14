package protocol

import "github.com/gin-gonic/gin"

type SetupFunc func(*Environment, *Middleware) ([]gin.HandlerFunc, ErrorCode)

type Endpoint struct {
	Method       Method
	RelativePath string
	Setup        SetupFunc
}
