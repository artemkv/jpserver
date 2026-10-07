package server

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func Serve(router *gin.Engine, port string) {
	ctx, restoreInterrupt := getNotifyContextForInterruptSignals()
	defer restoreInterrupt()

	httpServer := startServingAsync(router, port)
	waitForInterruptSignal(ctx)
	restoreInterrupt()
	shutDownWithTimeout(httpServer, 5*time.Second)
}

func getNotifyContextForInterruptSignals() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
}

func waitForInterruptSignal(ctx context.Context) {
	<-ctx.Done()
}

func startServingAsync(router *gin.Engine, port string) *http.Server {
	log.Printf("Starting server on port %s", port)
	httpServer := &http.Server{
		Addr:    port,
		Handler: router,
	}
	go listenAndServe(httpServer)
	return httpServer
}

func listenAndServe(httpServer *http.Server) {
	err := httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error serving: %s\n", err)
	}
}

func shutDownWithTimeout(httpServer *http.Server, timeout time.Duration) {
	log.Println("Shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
