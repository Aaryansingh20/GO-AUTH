// main.go - Update the session store setup

package main

import (
	"log"
	"os"

	controllers "github.com/Aaryansingh20/jwt/controllers"
	database "github.com/Aaryansingh20/jwt/database"
	models "github.com/Aaryansingh20/jwt/models"
	routes "github.com/Aaryansingh20/jwt/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth/gothic"
	"github.com/ulule/limiter/v3"
	ginLimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"
)

func main() {
	// Load .env
	if _, err := os.Stat(".env"); err == nil {
		godotenv.Load(".env")
		log.Println("✅ Loaded .env file")
	}

	// Setup session store for Goth (Google OAuth)
	key := os.Getenv("JWT_SECRET")
	if key == "" {
		key = "secret-key-please-change-this-in-production"
		log.Println("⚠️  Using default session key - please set JWT_SECRET in .env")
	}
	
	// Create session store with proper settings
	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(86400 * 30) // 30 days
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = false // Set to true in production with HTTPS
	store.Options.SameSite = 2   // SameSiteLaxMode
	
	gothic.Store = store
	log.Println("✅ Session store configured")

	// Initialize Google OAuth
	routes.SetupGoogleAuth()
	log.Println("✅ Google OAuth initialized")

	// Connect to database
	database.Client = database.DBinstance()
	database.Client.AutoMigrate(&models.User{})
	log.Println("✅ Database connected")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Rate limiting
	rate, err := limiter.NewRateFromFormatted("60-M")
	if err != nil {
		log.Fatal(err)
	}
	memoryStore := memory.NewStore()
	limiterMiddleware := ginLimiter.NewMiddleware(limiter.New(memoryStore, rate))

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(limiterMiddleware)

	// CORS Configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:3000",
		"http://localhost:8000",
	}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "token", "Authorization"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	router.Use(cors.New(config))

	// Regular auth routes
	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	// Google OAuth routes
	routes.GoogleAuthRoutes(router)

	// Google Auth Callback Handler (for frontend)
	router.GET("/api/auth/google/callback", controllers.GoogleAuthCallback())

	router.GET("/api-1", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-1"})
	})

	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-2"})
	})

	log.Printf("🚀 Server starting on port %s", port)
	log.Println("📍 Google OAuth: http://localhost:8000/auth/google")
	log.Println("📍 Test with: curl http://localhost:8000/auth/google")
	router.Run(":" + port)
}
