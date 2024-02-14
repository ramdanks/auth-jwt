package protocol

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const (
	ENV_MODE_DEVELOPMENT EnvironmentMode = 0
	ENV_MODE_STAGING     EnvironmentMode = 1
	ENV_MODE_PRODUCTION  EnvironmentMode = 2
)

type Environment struct {
	EnvironmentMode
}

func (e *Environment) Init() error {

	var envFilepath string
	var ginMode string

	switch e.EnvironmentMode {
	case ENV_MODE_DEVELOPMENT:
		ginMode = gin.DebugMode
		envFilepath = FILEPATH_ENV_DEVELOPMENT

	case ENV_MODE_STAGING:
		ginMode = gin.ReleaseMode
		envFilepath = FILEPATH_ENV_STAGING

	case ENV_MODE_PRODUCTION:
		ginMode = gin.ReleaseMode
		envFilepath = FILEPATH_ENV_PRODUCTION
	}

	gin.SetMode(ginMode)
	return godotenv.Load(envFilepath)
}

func (e *Environment) Lookup(key EnvKey) (string, bool) {
	return os.LookupEnv(key)
}

func (e *Environment) Get(key EnvKey) string {
	return os.Getenv(key)
}

type EnvironmentMode int
