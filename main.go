package main

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-rest-api/handlers"
)

func main() {

	router := gin.Default()

	customerHandler := handlers.NewCustomerHandler()
	customerHandler.SetupCustomerHandlers(router)

	router.Run(":8000")
}
