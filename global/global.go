package global

import (
	"fmt"
	"goAccounting/initialize"
	"log"
)

var (
	GlobalDb  = initialize.Database
	GlobalRdb = initialize.Rdb
	Config    = initialize.Config
	Cache     = initialize.Cache
)

func init() {
	fmt.Println("[global]: starting init")

	// Initialize resources
	GlobalDb = initialize.Database
	GlobalRdb = initialize.Rdb
	Config = initialize.Config
	Cache = initialize.Cache

	// Validate all resources are properly initialized
	if GlobalDb == nil {
		log.Fatal("global: failed to initialize database connection")
	}
	if GlobalRdb == nil {
		log.Fatal("global: failed to initialize Redis connection")
	}
	if Config == nil {
		log.Fatal("global: failed to load configuration")
	}
	if Cache == nil {
		log.Fatal("global: failed to initialize cache")
	}

	// Test database connectivity
	if err := testDatabaseConnection(); err != nil {
		log.Fatal("global: database connection test failed:", err)
	}

	// Test Redis connectivity
	if err := testRedisConnection(); err != nil {
		log.Fatal("global: Redis connection test failed:", err)
	}

	fmt.Println("[global]: init success - all services connected")
}

// testDatabaseConnection validates database connectivity
func testDatabaseConnection() error {
	if GlobalDb == nil {
		return fmt.Errorf("database instance is nil")
	}

	// Test database ping/connection
	sqlDB, err := GlobalDb.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %v", err)
	}

	fmt.Println("[global]: database connection verified")
	return nil
}

// testRedisConnection validates Redis connectivity
func testRedisConnection() error {
	if GlobalRdb == nil {
		return fmt.Errorf("Redis instance is nil")
	}

	// Test Redis ping
	if err := GlobalRdb.Ping(GlobalRdb.Context()).Err(); err != nil {
		return fmt.Errorf("Redis ping failed: %v", err)
	}

	fmt.Println("[global]: Redis connection verified")
	return nil
}

// HealthCheck provides a health check endpoint for frontend connectivity
func HealthCheck() map[string]string {
	status := make(map[string]string)

	// Check database
	if err := testDatabaseConnection(); err != nil {
		status["database"] = "unhealthy: " + err.Error()
	} else {
		status["database"] = "healthy"
	}

	// Check Redis
	if err := testRedisConnection(); err != nil {
		status["redis"] = "unhealthy: " + err.Error()
	} else {
		status["redis"] = "healthy"
	}

	// Check cache
	if Cache == nil {
		status["cache"] = "unhealthy: cache instance is nil"
	} else {
		status["cache"] = "healthy"
	}

	return status
}
