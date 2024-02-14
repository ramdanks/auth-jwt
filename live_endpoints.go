package main

import (
	"example/auth-jwt/endpoints"
	"example/auth-jwt/protocol"
)

var live_endpoints = [...]*protocol.Endpoint{
	endpoints.NewPing(),
	endpoints.NewLogin(),
	endpoints.NewRefresh(),
	endpoints.NewVerify(),
}
