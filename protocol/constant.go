package protocol

const (
	FILEPATH_ENV_DEVELOPMENT string = ".env/dev.env"
	FILEPATH_ENV_STAGING     string = ".env/stag.env"
	FILEPATH_ENV_PRODUCTION  string = ".env/prod.env"
)

const (
	ENV_KEY_ADDR             EnvKey = "ADDR"
	ENV_KEY_PORT             EnvKey = "PORT"
	ENV_KEY_ACCESS_TOKEN     EnvKey = "ACCESS_TOKEN_SECRET"
	ENV_KEY_REFRESH_TOKEN    EnvKey = "REFRESH_TOKEN_SECRET"
	ENV_KEY_ACCESS_EXP_TIME  EnvKey = "ACCESS_EXP_TIME"
	ENV_KEY_REFRESH_EXP_TIME EnvKey = "REFRESH_EXP_TIME"
)

const (
	ERR_NONE                               ErrorCode = 0
	ERR_SERVER_ENGINE_DETACHED             ErrorCode = 1
	ERR_ENV_INIT_FAILURE                   ErrorCode = 101
	ERR_ENV_ADDR_NOT_SPECIFIED             ErrorCode = 102
	ERR_ENV_PORT_NOT_SPECIFIED             ErrorCode = 103
	ERR_ENV_ACCESS_TOKEN_NOT_SPECIFIED     ErrorCode = 104
	ERR_ENV_REFRESH_TOKEN_NOT_SPECIFIED    ErrorCode = 105
	ERR_ENV_ACCESS_TOKEN_MALFORMED         ErrorCode = 106
	ERR_ENV_REFRESH_TOKEN_MALFORMED        ErrorCode = 107
	ERR_ENV_ACCESS_EXP_TIME_NOT_SPECIFIED  ErrorCode = 108
	ERR_ENV_ACCESS_EXP_TIME_MALFORMED      ErrorCode = 109
	ERR_ENV_REFRESH_EXP_TIME_NOT_SPECIFIED ErrorCode = 110
	ERR_ENV_REFRESH_EXP_TIME_MALFORMED     ErrorCode = 111
)

type EnvKey = string

type ErrorCode = int
