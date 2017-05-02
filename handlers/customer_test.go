package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-rest-api/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type CustomerDALMock struct {
	mock.Mock
}

func (d *CustomerDALMock) Get(id string) (models.Customer, error) {
	args := d.Called(id)
	return args.Get(0).(models.Customer), args.Error(1)
}

func (d *CustomerDALMock) GetAll() ([]models.Customer, error) {
	args := d.Called()
	return args.Get(0).([]models.Customer), args.Error(1)
}

func (d *CustomerDALMock) Create(customer models.Customer) error {
	args := d.Called()
	return args.Error(0)
}

func (d *CustomerDALMock) Update(id string, customer models.Customer) error {
	args := d.Called()
	return args.Error(0)
}

func NewCustomerHandlerMock(d *CustomerDALMock) *CustomerHandler {
	h := new(CustomerHandler)
	h.d = d
	return h
}

func getMockedRouter(d *CustomerDALMock) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	h := NewCustomerHandlerMock(d)
	h.SetupCustomerHandlers(router)
	return router
}

func TestCustomerHandler_GetAll(t *testing.T) {
	d := new(CustomerDALMock)
	d.On("GetAll").Return([]models.Customer{{Id: "123", Name: "John"}}, nil)
	r := getMockedRouter(d)

	req, _ := http.NewRequest("GET", "/customers", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, resp.Code, 200)

	var customersResp []models.Customer
	json.Unmarshal(resp.Body.Bytes(), &customersResp)

	assert.Len(t, customersResp, 1)
	assert.Equal(t, customersResp[0].Id, "123")
	assert.Equal(t, customersResp[0].Name, "John")
	assert.Len(t, customersResp[0].Phones, 0)
}

func TestCustomerHandler_GetAll_Error(t *testing.T) {
	d := new(CustomerDALMock)
	d.On("GetAll").Return([]models.Customer{{Id: "123", Name: "John"}}, errors.New("Test"))
	r := getMockedRouter(d)

	req, _ := http.NewRequest("GET", "/customers", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, resp.Code, 500)

	var errorResp gin.H
	json.Unmarshal(resp.Body.Bytes(), &errorResp)
	assert.Equal(t, errorResp["error"], "Could not get Customers list")
}

func TestCustomerHandler_Get(t *testing.T) {
	d := new(CustomerDALMock)
	d.On("Get", "123").Return(models.Customer{Id: "123", Name: "John"}, nil)
	r := getMockedRouter(d)

	req, _ := http.NewRequest("GET", "/customers/123", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, resp.Code, 200)

	var customerResp models.Customer
	json.Unmarshal(resp.Body.Bytes(), &customerResp)

	assert.Equal(t, customerResp.Id, "123")
	assert.Equal(t, customerResp.Name, "John")
	assert.Len(t, customerResp.Phones, 0)
}

func TestCustomerHandler_Get_Error(t *testing.T) {
	d := new(CustomerDALMock)
	d.On("Get", "123").Return(models.Customer{Id: "123", Name: "John"}, errors.New("Test"))
	r := getMockedRouter(d)

	req, _ := http.NewRequest("GET", "/customers/123", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, resp.Code, 404)

	var errorResp gin.H
	json.Unmarshal(resp.Body.Bytes(), &errorResp)
	assert.Equal(t, errorResp["error"], "Customer not found")
}

func TestCustomerHandler_Create(t *testing.T) {
	d := new(CustomerDALMock)
	d.On("Create").Return(nil)
	r := getMockedRouter(d)

	newCustomer := models.Customer{
		Id:   "123",
		Name: "John",
	}

	newCustomerJson, _ := json.Marshal(newCustomer)
	req, _ := http.NewRequest("POST", "/customers", bytes.NewBuffer(newCustomerJson))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, 201, resp.Code)

	var successResp map[string]models.Customer
	json.Unmarshal(resp.Body.Bytes(), &successResp)
	assert.Equal(t, newCustomer, successResp["success"])
}

func TestCustomerHandler_Create_Error(t *testing.T) {
	d := new(CustomerDALMock)
	d.On("Create").Return(errors.New("Test"))
	r := getMockedRouter(d)

	newCustomer := models.Customer{
		Id:   "123",
		Name: "John",
	}

	newCustomerJson, _ := json.Marshal(newCustomer)
	req, _ := http.NewRequest("POST", "/customers", bytes.NewBuffer(newCustomerJson))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, 500, resp.Code)

	var errorResp gin.H
	json.Unmarshal(resp.Body.Bytes(), &errorResp)
	assert.Equal(t, "Could not create customer (John)", errorResp["error"])
}

func TestCustomerHandler_Create_InvalidPayload(t *testing.T) {
	d := new(CustomerDALMock)
	d.On("Create").Return(nil)
	r := getMockedRouter(d)

	newCustomer := models.Customer{
		Id:   "123",
		Name: "",
	}

	newCustomerJson, _ := json.Marshal(newCustomer)
	req, _ := http.NewRequest("POST", "/customers", bytes.NewBuffer(newCustomerJson))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, 422, resp.Code)

	var errorResp gin.H
	json.Unmarshal(resp.Body.Bytes(), &errorResp)
	assert.Equal(t, errorResp["error"], "Customer name could not be empty")
}

func TestCustomerHandler_Update(t *testing.T) {
	d := new(CustomerDALMock)
	d.On("Get", "123").Return(models.Customer{Id: "123", Name: "John"}, nil)
	d.On("Update").Return(nil)
	r := getMockedRouter(d)

	updatedCustomer := models.Customer{
		Id:   "123",
		Name: "John Test",
	}

	updatedCustomerJson, _ := json.Marshal(updatedCustomer)
	req, _ := http.NewRequest("PUT", "/customers/123", bytes.NewBuffer(updatedCustomerJson))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)

	var successResp map[string]models.Customer
	json.Unmarshal(resp.Body.Bytes(), &successResp)
	assert.Equal(t, updatedCustomer, successResp["success"])
}

func TestCustomerHandler_Update_InvalidCustomer(t *testing.T) {
	d := new(CustomerDALMock)
	d.On("Get", "123").Return(models.Customer{}, errors.New("Test"))
	r := getMockedRouter(d)

	updatedCustomer := models.Customer{
		Id:   "123",
		Name: "John Test",
	}

	updatedCustomerJson, _ := json.Marshal(updatedCustomer)
	req, _ := http.NewRequest("PUT", "/customers/123", bytes.NewBuffer(updatedCustomerJson))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, 404, resp.Code)

	var errorResp gin.H
	json.Unmarshal(resp.Body.Bytes(), &errorResp)
	assert.Equal(t, "Customer not found", errorResp["error"])
}

func TestCustomerHandler_Update_InvalidCustomerID(t *testing.T) {
	d := new(CustomerDALMock)
	d.On("Get", "123").Return(models.Customer{Id: "123"}, nil)
	r := getMockedRouter(d)

	updatedCustomer := models.Customer{
		Id:   "1234",
		Name: "John Test",
	}

	updatedCustomerJson, _ := json.Marshal(updatedCustomer)
	req, _ := http.NewRequest("PUT", "/customers/123", bytes.NewBuffer(updatedCustomerJson))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, 422, resp.Code)

	var errorResp gin.H
	json.Unmarshal(resp.Body.Bytes(), &errorResp)
	assert.Equal(t, "Could not update customer (Invalid ID)", errorResp["error"])
}

func TestCustomerHandler_Update_Error(t *testing.T) {
	d := new(CustomerDALMock)
	d.On("Get", "123").Return(models.Customer{Id: "123"}, nil)
	d.On("Update").Return(errors.New("Test"))
	r := getMockedRouter(d)

	updatedCustomer := models.Customer{
		Id:   "123",
		Name: "John Test",
	}

	updatedCustomerJson, _ := json.Marshal(updatedCustomer)
	req, _ := http.NewRequest("PUT", "/customers/123", bytes.NewBuffer(updatedCustomerJson))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, 422, resp.Code)

	var errorResp gin.H
	json.Unmarshal(resp.Body.Bytes(), &errorResp)
	assert.Equal(t, "Could not update customer", errorResp["error"])
}
