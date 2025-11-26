// routes/googleAuthRoutes.go
package routes

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func SetupGoogleAuth() {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	callbackURL := os.Getenv("GOOGLE_CALLBACK_URL")

	if clientID == "" || clientSecret == "" {
		log.Fatal("GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET must be set")
	}

	log.Printf("Google OAuth Client ID: %s", clientID[:20]+"...")
	log.Printf("Google OAuth Callback URL: %s", callbackURL)

	goth.UseProviders(
		google.New(clientID, clientSecret, callbackURL, "email", "profile"),
	)
}

func GoogleAuthRoutes(router *gin.Engine) {
	// Start Google OAuth flow
	router.GET("/auth/google", func(c *gin.Context) {
		log.Println("Starting Google OAuth flow...")
		
		// Set provider in query
		q := c.Request.URL.Query()
		q.Add("provider", "google")
		c.Request.URL.RawQuery = q.Encode()

		// Start OAuth
		gothic.BeginAuthHandler(c.Writer, c.Request)
	})

	// Google callback
	router.GET("/auth/google/callback", func(c *gin.Context) {
		log.Println("Received Google callback...")
		
		// Set provider in query
		q := c.Request.URL.Query()
		q.Add("provider", "google")
		c.Request.URL.RawQuery = q.Encode()

		// Complete authentication
		user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
		if err != nil {
			log.Printf("❌ Error in Google callback: %v", err)
			
			// Try to get the error details
			if err.Error() == "gothic: user data not found" {
				log.Println("Session error - trying alternative approach...")
				c.Redirect(302, "http://localhost:3000/signup?error=google_auth_failed")
				return
			}
			
			c.JSON(500, gin.H{"error": "Authentication failed", "details": err.Error()})
			return
		}

		log.Printf("✅ Google user authenticated: %s (%s)", user.Name, user.Email)
		log.Printf("User ID: %s", user.UserID)
		log.Printf("Avatar: %s", user.AvatarURL)

		// Redirect to frontend callback with user info
		frontendURL := "http://localhost:3000/auth/callback"
		redirectURL := frontendURL + "?email=" + user.Email + "&name=" + user.Name + "&picture=" + user.AvatarURL + "&provider=google"
		
		log.Printf("Redirecting to: %s", redirectURL)
		c.Redirect(302, redirectURL)
	})
}
