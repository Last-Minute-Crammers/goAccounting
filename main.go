package main

import (
	"net/http"

	"context"
	"log"
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
