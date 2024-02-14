package main

import (
	"encoding/hex"
	"example/auth-jwt/protocol"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

const (
	ENV_MODE = protocol.ENV_MODE_DEVELOPMENT
)

func main() {

	env := protocol.Environment{EnvironmentMode: ENV_MODE}

	err := env.Init()
	if err != nil {
		os.Exit(protocol.ERR_ENV_INIT_FAILURE)
		return
	}

	v, ok := env.Lookup(protocol.ENV_KEY_ACCESS_TOKEN)
	if !ok {
		os.Exit(protocol.ERR_ENV_ACCESS_TOKEN_NOT_SPECIFIED)
		return
	}

	token, err := hex.DecodeString(v)
	if err != nil {
		os.Exit(protocol.ERR_ENV_ACCESS_TOKEN_MALFORMED)
		return
	}

	m := protocol.Middleware{AccessTokenSecret: token}

	engine := gin.Default()

	for _, ep := range live_endpoints {
		handlers, errcode := ep.Setup(&env, &m)
		if errcode != protocol.ERR_NONE {
			os.Exit(errcode)
			return
		}

		engine.Handle(ep.Method, ep.RelativePath, handlers...)
	}

	port, ok := env.Lookup(protocol.ENV_KEY_PORT)
	if !ok {
		os.Exit(protocol.ERR_ENV_PORT_NOT_SPECIFIED)
	}

	addr, ok := env.Lookup(protocol.ENV_KEY_ADDR)
	if !ok {
		os.Exit(protocol.ERR_ENV_ADDR_NOT_SPECIFIED)
	}

	path := fmt.Sprintf("%s:%s", addr, port)
	err = engine.Run(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(protocol.ERR_SERVER_ENGINE_DETACHED)
	}

}
