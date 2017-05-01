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

func (d CustomerDALMock) Get(id string) (models.Customer, error) {
	args := d.Called(id)
	return args.Get(0).(models.Customer), args.Error(1)
}

func (d CustomerDALMock) GetAll() ([]models.Customer, error) {
	args := d.Called()
	return args.Get(0).([]models.Customer), args.Error(1)
}

func (d CustomerDALMock) Create(customer models.Customer) error {
	args := d.Called()
	return args.Error(0)
}

func (d CustomerDALMock) Update(id string, customer models.Customer) error {
	args := d.Called()
	return args.Error(0)
}

func NewCustomerHandlerMock(d *CustomerDALMock) *CustomerHandler {
	h := new(CustomerHandler)
	h.d = d
	return h
}

func getMockRouter(h *CustomerHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h.SetupCustomerHandlers(router)
	return router
}

func TestCustomerHandler_GetAll(t *testing.T) {
	d := new(CustomerDALMock)
	d.On("GetAll").Return([]models.Customer{{Id: "123", Name: "John"}}, nil)
	h := NewCustomerHandlerMock(d)

	r := getMockRouter(h)

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
	h := NewCustomerHandlerMock(d)

	r := getMockRouter(h)

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
	h := NewCustomerHandlerMock(d)

	r := getMockRouter(h)

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
	h := NewCustomerHandlerMock(d)

	r := getMockRouter(h)

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
	h := NewCustomerHandlerMock(d)

	r := getMockRouter(h)

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
}

func TestCustomerHandler_Create_Error(t *testing.T) {
	d := new(CustomerDALMock)
	d.On("Create").Return(errors.New("Test"))
	h := NewCustomerHandlerMock(d)

	r := getMockRouter(h)

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
	assert.Equal(t, errorResp["error"], "Could not create customer (John)")
}

func TestCustomerHandler_Create_InvalidPayload(t *testing.T) {
	d := new(CustomerDALMock)
	d.On("Create").Return(nil)
	h := NewCustomerHandlerMock(d)

	r := getMockRouter(h)

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
