package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func DBinstance() *gorm.DB {
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatal("Error loading the .env file")
    }

    // Get PostgreSQL URL from .env (instead of MongoDB)
    PostgresDB := os.Getenv("DATABASE_URL")

    // Connect to PostgreSQL (instead of MongoDB)
    client, err := gorm.Open(postgres.Open(PostgresDB), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Connected to PostgreSQL!!")
    return client
}

var Client *gorm.DB = DBinstance()

// Note: OpenCollection is not needed for PostgreSQL/GORM
// GORM works directly with models, no need for collections
