package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"email-verifier/internal/handler"
)

func main() {
	// --------------------------------------------------
	// Gin production mode
	// --------------------------------------------------
	gin.SetMode(gin.ReleaseMode)

	// --------------------------------------------------
	// Create Gin engine (NO default logger spam)
	// --------------------------------------------------
	r := gin.New()
	r.Use(gin.Recovery())

	// --------------------------------------------------
	// Basic CORS (safe for public utility)
	// --------------------------------------------------
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// --------------------------------------------------
	// Serve UI (Plain HTML & CSS)
	// --------------------------------------------------
	r.StaticFile("/", "./web/index.html")
	r.Static("/styles.css", "./web/styles.css")

	// --------------------------------------------------
	// API Routes
	// --------------------------------------------------
	r.POST("/verify", handler.VerifyHandler)

	// --------------------------------------------------
	// Port (required for Render / Fly / Railway)
	// --------------------------------------------------
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	// --------------------------------------------------
	// Start server
	// --------------------------------------------------
	r.Run(":" + port)
}
