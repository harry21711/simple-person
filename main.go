package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomepageHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Simple Person listing API"})
}
func main() {
	fmt.Println("Hi")

	router := gin.Default()
	router.GET("/", HomepageHandler)
	router.Run("localhost:8083")
}
