package main

import (
	"artemkv.net/jpserver/app"
	"artemkv.net/jpserver/server"
	"github.com/gin-gonic/gin"
)

func main() {
	// load .env
	LoadDotEnv()

	// determine port
	port := GetOptionalString("JPSERVER_PORT", ":9999")

	// configure router
	allowedOrigins := GetOptionalString(
		"JPSERVER_ALLOW_ORIGIN", "https://localhost")
	router := gin.New()
	app.SetupRouter(router, allowedOrigins)

	// start the server
	server.Serve(router, port)
}
