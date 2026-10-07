package app

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine, allowedOrigins string) {
	// setup logger
	router.Use(requestLogger())

	// setup CORS
	allowedOriginsArr := strings.Split(allowedOrigins, ",")
	router.Use(cors.New(getCorsConfig(allowedOriginsArr)))

	// do business
	router.GET("/frame", handleFrame)

	// handle 404
	router.NoRoute(notFoundHandler())
}

func getCorsConfig(allowedOrigins []string) cors.Config {
	return cors.Config{
		AllowOrigins:  allowedOrigins,
		AllowHeaders:  []string{"*"},
		AllowMethods:  []string{"*"},
		ExposeHeaders: []string{"*"},
	}
}

func toBinary(c *gin.Context, data []byte) {
	c.Data(http.StatusOK, "application/octet-stream", data)
}

func toBadRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
}

func toInternalServerError(c *gin.Context, errText string) {
	c.JSON(http.StatusInternalServerError, gin.H{"err": errText})
}

func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		message := fmt.Sprintf("%s %s %s (Origin: %s)",
			c.ClientIP(),
			c.Request.Method,
			c.Request.URL.Path,
			c.Request.Header.Get("Origin"),
		)
		log.Print(message)
	}
}

func notFoundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"err": "Not found"})
	}
}
