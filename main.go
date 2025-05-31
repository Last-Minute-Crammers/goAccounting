package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "goAccounting/global"
	"goAccounting/initialize"
	"goAccounting/router"
)

var httpServer *http.Server

func main() {
	_ = initialize.Config

	router.RegisterPublicRoutes()
	router.RegisterAIRoutes()

	// 打印路由表
	foundPingRoute := false
	for _, r := range router.Engine.Routes() {
		log.Printf("Registered route: %s %s", r.Method, r.Path)
		if r.Path == "/api/public/ping" {
			foundPingRoute = true
		}
	}

	if !foundPingRoute {
		log.Fatal("Route /api/public/ping not found")
	}

	router.Engine.Run(":8080")
}

func shutDown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server ...")
	if err := httpServer.Shutdown(context.TODO()); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	log.Println("Server exiting")
}
