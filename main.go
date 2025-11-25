package main

import (
	"log"
	"os"

	database "github.com/Aaryansingh20/jwt/database"
	models "github.com/Aaryansingh20/jwt/models"
	routes "github.com/Aaryansingh20/jwt/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ulule/limiter/v3"
	ginLimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"
)

func main() {
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatal("Error loading the .env file")
    }

    // Connect to PostgreSQL database
    database.Client = database.DBinstance()

    // Auto-migrate (create tables automatically)
    database.Client.AutoMigrate(&models.User{})

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    // Rate limiting setup (60 requests per minute per IP)
    rate, err := limiter.NewRateFromFormatted("60-M")
    if err != nil {
        log.Fatal(err)
    }
    store := memory.NewStore()
    limiterMiddleware := ginLimiter.NewMiddleware(limiter.New(store, rate))

    router := gin.New()
    router.Use(gin.Logger())
    router.Use(limiterMiddleware) // <-- Add rate limiting globally

    // Add CORS for Next.js frontend
    config := cors.DefaultConfig()
    config.AllowOrigins = []string{"http://localhost:3000"} // Your Next.js URL
    config.AllowCredentials = true
    config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "token"}
    config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
    router.Use(cors.New(config))

    routes.AuthRoutes(router)
    routes.UserRoutes(router)

    router.GET("/api-1", func(c *gin.Context) {
        c.JSON(200, gin.H{"success": "Access granted for api-1"})
    })

    router.GET("/api-2", func(c *gin.Context) {
        c.JSON(200, gin.H{"success": "Access granted for api-2"})
    })

    log.Printf("🚀 Server starting on port %s", port)
    router.Run(":" + port)
}
