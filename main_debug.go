package main

import (
    "fmt"
    "log"
    "net/http"
    "time"

    // 暂时注释掉 global 导入
    // _ "goAccounting/global"
    "goAccounting/initialize"
    "goAccounting/router"
)

func main() {
    log.Println("=== MAIN FUNCTION STARTED ===")
    
    // 使用 initialize.go 中的 Config
    config := initialize.Config
    log.Printf("Config type: %T", config)
    log.Printf("Loaded config: %+v", config)
    
    // 获取端口，添加默认值处理
    port := 8080  // 默认端口
    if config != nil && config.System.Addr != 0 {
        port = config.System.Addr
    }
    
    // 启动HTTP服务器
    addr := fmt.Sprintf(":%d", port)
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
    
    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("Failed to start server: %v", err)
    }
}