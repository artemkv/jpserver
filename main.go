package main

import (
	"crypto/rand"
	"fmt"
	"io"
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

	// control access
	accessCode := GetOptionalString("JPSERVER_CODE", "")
	if accessCode == "" {
		var err error
		accessCode, err = generateAccessCode()
		if err != nil {
			log.Fatal("Could not generate the access code")
		}
	}
	fmt.Printf("***** ENDPOINT: %s, ACCESS CODE: %s *****\n", endpoint, accessCode)

	// configure router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	app.SetupRouter(router, accessCode)

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

func generateAccessCode() (string, error) {
	const table = "0123456789"
	buffer := make([]byte, 6)
	if _, err := io.ReadFull(rand.Reader, buffer); err != nil {
		return "", err
	}
	for i, b := range buffer {
		buffer[i] = table[b%byte(len(table))]
	}
	return string(buffer), nil
}
