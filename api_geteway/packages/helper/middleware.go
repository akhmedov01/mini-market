package helper

import (
	"fmt"
	"main/config"
	"main/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func StartMiddleware(c *gin.Context) {
	// before request
	fmt.Printf("%s request start path: %s time %v\n", c.Request.Method, c.Request.URL.Path, time.Now())
	c.Next()
}
func EndMiddleware(c *gin.Context) {
	c.Next()
	// after request
	fmt.Printf("%s request end path: %s time %v\n", c.Request.Method, c.Request.URL.Path, time.Now())
}
func LoggerAllInOne(c *gin.Context) {
	// before request
	t := time.Now()
	c.Next()
	// after request
	latency := time.Since(t)

	// access the status we are sending
	status := c.Writer.Status()

	fmt.Printf("Completed %s %s with status code %d in %v\n", c.Request.Method, c.Request.URL.Path, status, latency)
}
func Logger(c *gin.Context) {
	// before request
	beforeRequest(c)
	c.Next()
	// after request
	afterRequest(c)
}
func beforeRequest(c *gin.Context) {
	// before request
	t := time.Now()
	c.Set("start", t)
	c.Next()
}
func afterRequest(c *gin.Context) {
	// Get the start time from the request context
	startTime, exists := c.Get("start")
	if !exists {
		startTime = time.Now()
	}

	// Calculate the request duration
	duration := time.Since(startTime.(time.Time))

	// Log the request completion time and duration
	fmt.Printf("Completed %s %s in %v\n", c.Request.Method, c.Request.URL.Path, duration)
}

// AuthMiddleWare is a middleware function for authentication
func AuthMiddleWare(c *gin.Context) {
	token := c.GetHeader("Authorization")

	if token == "" {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    "UNAUTHORIZED!",
			"message": "Token not found...",
		})
		return
	}

	staffInfo, err := ParseClaims(token, config.JWTSecretKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    "INVALID TOKEN!",
			"message": "Provided token is not valid...",
		})
		return
	}

	c.Set("staff_info", staffInfo)
	c.Next()
}

func ValidatePasswordMiddleware(c *gin.Context) {
	var requestBody struct {
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		c.Abort()
		return
	}

	password := requestBody.Password

	if !IsValidPassword(password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password "})
		c.Abort()
		return
	}

	c.Next()
}

func ValidateUsernameMiddleware(c *gin.Context) {
	var req models.LoginReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		c.Abort()
		return
	}

	if !IsValidLogin(req.Loging) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username"})
		c.Abort()
		return
	}

	c.Next()
}
