package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"email-verifier/internal/model"
	"email-verifier/internal/service"
)

// VerifyHandler handles POST /verify
func VerifyHandler(c *gin.Context) {
	var req model.VerifyRequest

	// ------------------------------
	// Parse JSON body
	// ------------------------------
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON body",
		})
		return
	}

	// ------------------------------
	// Basic validation
	// ------------------------------
	if req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email is required",
		})
		return
	}

	if len(req.Email) > 254 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email length",
		})
		return
	}

	// ------------------------------
	// Call core service
	// ------------------------------
	result := service.VerifyEmail(req.Email)

	// ------------------------------
	// Return result
	// ------------------------------
	c.JSON(http.StatusOK, result)
}
