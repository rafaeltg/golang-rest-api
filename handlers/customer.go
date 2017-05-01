package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-rest-api/dal"
	"github.com/golang-rest-api/models"
)

type CustomerHandler struct {
	d dal.BaseDAL
}

func NewCustomerHandler() *CustomerHandler {
	h := new(CustomerHandler)
	h.d = dal.NewCustomerDAL()
	return h
}

func (h CustomerHandler) SetupCustomerHandlers(router *gin.Engine) {
	router.GET("/customers", h.GetAll)
	router.GET("/customers/:id", h.Get)
	router.POST("/customers", h.Create)
	router.PUT("/customers/:id", h.Update)
}

func (h CustomerHandler) Get(c *gin.Context) {
	customerId := c.Params.ByName("id")
	customer, err := h.d.Get(customerId)
	if err != nil {
		c.JSON(404, gin.H{"error": "Customer not found"})
	} else {
		c.JSON(200, customer)
	}
}

func (h CustomerHandler) GetAll(c *gin.Context) {
	customers, err := h.d.GetAll()
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not get Customers list"})
	} else {
		c.JSON(200, customers)
	}
}

func (h CustomerHandler) Create(c *gin.Context) {
	var customerJson models.Customer
	c.Bind(&customerJson)

	fmt.Println(customerJson)

	if customerJson.Name == "" {
		fmt.Println("ERRRR")
		c.JSON(422, gin.H{"error": "Customer name could be empty"})
		return
	}

	fmt.Println("AAAAAAAAAAAAA", customerJson)
	if err := h.d.Create(customerJson); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Could not create customer (%v)", customerJson)})
	} else {
		c.JSON(201, gin.H{"success": customerJson})
	}
}

func (h CustomerHandler) Update(c *gin.Context) {
	customerId := c.Params.ByName("id")
	customer, err := h.d.Get(customerId)

	if customer.Id == "" || err != nil {
		c.JSON(404, gin.H{"error": "Customer not found"})
		return
	}

	var newCustomer models.Customer
	c.Bind(&newCustomer)

	if customerId != newCustomer.Id {
		c.JSON(422, gin.H{"error": "Could not update customer (Invalid ID)"})
	}

	if err = h.d.Update(customer.Id, newCustomer); err != nil {
		c.JSON(422, gin.H{"error": "Could not update customer"})
	} else {
		c.JSON(200, gin.H{"success": newCustomer})
	}
}
