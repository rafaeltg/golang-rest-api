package main

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-rest-api/handlers"
	"os"
	"fmt"
)

func main() {

	router := gin.Default()

	// Init Handlers
	customerHandler := handlers.NewCustomerHandler()
	customerHandler.SetupCustomerHandlers(router)

	// Get server port
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8000"
	}

	router.Run(fmt.Sprintf(":%s", port))
}
