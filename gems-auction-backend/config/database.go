package config

import (
	"context"
	"fmt"
	"log"
	"net/url" 
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

// ConnectDatabase initializes PostgreSQL connection pool
func ConnectDatabase() {

	encodedPassword := url.QueryEscape(AppConfig.DBPassword)

	dsn := fmt.Sprintf(
	"postgres://%s:%s@%s:%s/%s",
	AppConfig.DBUser,
	encodedPassword,
	AppConfig.DBHost,
	AppConfig.DBPort,
	AppConfig.DBName,
)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal("❌ Unable to parse DB config:", err)
	}

	config.MaxConns = AppConfig.MaxDBConns

	DB, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal("❌ Unable to create DB pool:", err)
	}

	// Check DB connection
	err = DB.Ping(ctx)
	if err != nil {
		log.Fatal("❌ Database not reachable:", err)
	}

	log.Println("✅ PostgreSQL Connected Successfully")
}

// CloseDatabase gracefully closes DB connection
func CloseDatabase() {
	if DB != nil {
		DB.Close()
		log.Println("🛑 Database connection closed")
	}
}
