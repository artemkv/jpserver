package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"artemkv.net/jpserver/app"
	"artemkv.net/jpserver/server"
	"github.com/gin-gonic/gin"
	"github.com/mdp/qrterminal/v4"
)

const (
	ACCESS_CODE_LENGTH = 8
	MIN_PORT           = 9991
	MAX_PORT           = 9999
)

func main() {
	LoadDotEnv()

	host := GetOptionalString("JPSERVER_HOST", getDefaultHost)
	port := GetOptionalString("JPSERVER_PORT",
		func() string { return getDefaultPort(host) })
	accessCode := GetOptionalString("JPSERVER_CODE", getDefaultAccessCode)

	addr := fmt.Sprintf("%s:%s", host, port)
	showConnectionInfo(addr, accessCode)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	app.SetupRouter(router, accessCode)
	server.Serve(router, addr)
}

func getDefaultHost() string {
	host, err := detectLocalIP()
	if err != nil {
		log.Fatal(err)
	}
	return host
}

func getDefaultPort(host string) string {
	port, err := getFreePort(host, MIN_PORT, MAX_PORT)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%d", port)
}

func getDefaultAccessCode() string {
	accessCode, err := generateAccessCode(ACCESS_CODE_LENGTH)
	if err != nil {
		log.Fatal(err)
	}
	return accessCode
}

func detectLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", fmt.Errorf("failed to detect local IP: %w", err)
	}
	defer conn.Close()
	localIP := conn.LocalAddr().(*net.UDPAddr).IP
	return localIP.String(), nil
}

func getFreePort(host string, minPort int, maxPort int) (int, error) {
	for port := minPort; port <= maxPort; port++ {
		addr := fmt.Sprintf("%s:%d", host, port)
		ln, err := net.Listen("tcp", addr)
		if err == nil {
			_ = ln.Close()
			return port, nil
		}
	}
	return 0, fmt.Errorf("failed to find a free port for %s", host)
}

func generateAccessCode(length int) (string, error) {
	const table = "0123456789"
	buffer := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, buffer); err != nil {
		return "", fmt.Errorf("failed to generate access code: %w", err)
	}
	for i, b := range buffer {
		buffer[i] = table[b%byte(len(table))]
	}
	return string(buffer), nil
}

func showConnectionInfo(addr string, accessCode string) {
	fmt.Printf("***** ADDRESS: %s, ACCESS CODE: %s *****\n", addr, accessCode)
	config := qrterminal.Config{
		Level:     qrterminal.M,
		Writer:    os.Stdout,
		BlackChar: "  ",
		WhiteChar: "██",
		QuietZone: 1,
	}
	qrterminal.GenerateWithConfig(fmt.Sprintf("%s|%s", addr, accessCode), config)
}
