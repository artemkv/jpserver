package main

import (
	"fmt"
	"net"
	"strings"

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
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	app.SetupRouter(router, allowedOrigins)

	// helps to connect
	suggestConnectionString(port)

	// start the server
	server.Serve(router, port)
}

func suggestConnectionString(port string) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			ip := ipNet.IP

			if ip.IsLoopback() || ip.To4() == nil {
				continue
			}

			if strings.HasPrefix(ip.String(), "192.168.0.") {
				fmt.Printf("***** Use this endpoint: %s:%s *****\n", ip.String(), port)
			}
		}
	}
}
