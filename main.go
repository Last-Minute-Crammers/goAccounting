package main

import (
	"fmt"
	"net/http"
	"time"

	//"fmt"
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

	httpServer = &http.Server{
		Addr:           fmt.Sprintf(":%d", initialize.Config.System.Addr),
		Handler:        router.Engine,
		WriteTimeout:   5 * time.Second,
		ReadTimeout:    5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := httpServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
	shutDown()
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
