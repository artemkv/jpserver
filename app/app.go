package app

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine, accessCode string) {
	// setup logger
	router.Use(requestLogger())

	// setup CORS
	router.Use(cors.New(getCorsConfig()))

	// do business
	router.GET("/frame", HandleWithAccessCode(accessCode, handleFrame))

	// handle 404
	router.NoRoute(notFoundHandler())
}

func getCorsConfig() cors.Config {
	return cors.Config{
		AllowOrigins:  []string{"*"},
		AllowHeaders:  []string{"*"},
		AllowMethods:  []string{"*"},
		ExposeHeaders: []string{"*"},
	}
}

func HandleWithAccessCode(accessCode string, handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Query("code")
		if code != accessCode {
			c.Status(http.StatusUnauthorized)
			return
		}
		handler(c)
	}
}

func toBinary(c *gin.Context, data []byte) {
	c.Data(http.StatusOK, "application/octet-stream", data)
}

func toBinaryGzip(c *gin.Context, data []byte) {
	c.Header("Content-Encoding", "gzip")
	c.Header("Vary", "Accept-Encoding")
	c.Data(http.StatusOK, "application/octet-stream", data)
}

func gzipBytes(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)

	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}

	// Must close the writer to flush gzip data
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func toBadRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
}

func toInternalServerError(c *gin.Context, errText string) {
	c.JSON(http.StatusInternalServerError, gin.H{"err": errText})
}

func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		message := fmt.Sprintf("%s %s %s",
			c.ClientIP(),
			c.Request.Method,
			c.Request.URL.String(),
		)
		log.Print(message)
	}
}

func notFoundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"err": "Not found"})
	}
}

func getQueryParamAsInt(c *gin.Context, key string, defaultValue int) int {
	val := c.Query(key)
	n, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return n
}
