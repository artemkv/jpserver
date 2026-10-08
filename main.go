package main

import (
	"fmt"
	"log"
	"net"

	"artemkv.net/jpserver/app"
	"artemkv.net/jpserver/server"
	"github.com/gin-gonic/gin"
)

func main() {
	// load .env
	LoadDotEnv()

	// connection endpoint
	host := GetOptionalString("JPSERVER_HOST", detectLocalIP())
	port := GetOptionalString("JPSERVER_PORT", "9999")
	endpoint := fmt.Sprintf("%s:%s", host, port)

	// configure router
	allowedOrigins := GetOptionalString(
		"JPSERVER_ALLOW_ORIGIN", "https://localhost")
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	app.SetupRouter(router, allowedOrigins)

	// start the server
	server.Serve(router, endpoint)
}

func detectLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Print("Could not detect local IP")
		return ""
	}
	defer conn.Close()
	localIP := conn.LocalAddr().(*net.UDPAddr).IP
	return localIP.String()
}
