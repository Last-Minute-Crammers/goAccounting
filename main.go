package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "goAccounting/global"
	"goAccounting/initialize"
	"goAccounting/router"
)

func main() {
	// 初始化配置
	config := initialize.Config
	
	// 启动HTTP服务器
	addr := fmt.Sprintf(":%d", config.System.Addr)
	log.Printf("Starting server on %s", addr)
	
	server := &http.Server{
		Addr:           addr,
		Handler:        router.Engine,
		WriteTimeout:   30 * time.Second,
		ReadTimeout:    30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	
	log.Printf("Server is running on http://localhost%s", addr)
	log.Printf("Health check: http://localhost%s/health", addr)
	log.Printf("Category API: http://localhost%s/api/v1/user/category/list", addr)
	
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
