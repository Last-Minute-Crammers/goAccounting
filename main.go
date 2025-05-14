package main

import (
	"net/http"
	//"fmt"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var httpServer *http.Server

func main() {
	httpServer = &http.Server{
		// Addr: ,
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
	log.Println("Server exiti程度g")
}
