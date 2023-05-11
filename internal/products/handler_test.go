package products

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createServer(service Service) *gin.Engine {
	handler := NewHandler(service)
	r := gin.Default()
	rg := r.Group("/api/v1")
	{
		prodRoute := rg.Group("/products")
		{
			prodRoute.GET("", handler.GetProducts)
		}
	}
	return r
}

func createRequestTest(method string, url string, body string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))

	req.Header.Add("Content-Type", "application/json")
	return req, httptest.NewRecorder()
}

type MockService struct {
	mock.Mock
}

func NewMockService() MockService {
	return MockService{}
}

func (s *MockService) GetAllBySeller(sellerID string) ([]Product, error) {
	args := s.Mock.Called(sellerID)
	pl := args.Get(0).([]Product)
	err := args.Error(1)

	return pl, err
}

func TestGetProduct(t *testing.T) {
	t.Run("Sucess", func(t *testing.T) {
		//Arrange
		service := NewMockService()
		server := createServer(&service)

		sellerId := "1"

		service.On("GetAllBySeller", sellerId).Return([]Product{
			{
				ID:          "1",
				SellerID:    "1",
				Description: "The description",
				Price:       100,
			},
			{
				ID:          "2",
				SellerID:    "1",
				Description: "The description",
				Price:       200,
			},
		}, nil)

		expectedHeaders := http.Header{
			"Content-Type": []string{"application/json; charset=utf-8"},
		}
		expectedBody := `[
			{
				"ID":          "1",
				"SellerID":    "1",
				"Description": "The description",
				"Price":       100
			},
			{
				"ID":          "2",
				"SellerID":    "1",
				"Description": "The description",
				"Price":       200
			}
		]`

		//Act
		request, respRecorder := createRequestTest(http.MethodGet, "/api/v1/products?seller_id=1", ``)
		server.ServeHTTP(respRecorder, request)

		//Assert
		assert.Equal(t, http.StatusOK, respRecorder.Code)
		assert.Equal(t, expectedHeaders, respRecorder.Header())
		assert.JSONEq(t, expectedBody, respRecorder.Body.String())
		service.AssertExpectations(t)
	})

	t.Run("SellerID Required", func(t *testing.T) {
		//Arrange
		service := NewMockService()
		server := createServer(&service)

		expectedHeaders := http.Header{
			"Content-Type": []string{"application/json; charset=utf-8"},
		}
		expectedBody := `{
			"error": "seller_id query param is required"
		}`

		//Act
		request, respRecorder := createRequestTest(http.MethodGet, "/api/v1/products", ``)
		server.ServeHTTP(respRecorder, request)

		//Assert
		assert.Equal(t, http.StatusBadRequest, respRecorder.Code)
		assert.Equal(t, expectedHeaders, respRecorder.Header())
		assert.JSONEq(t, expectedBody, respRecorder.Body.String())
	})

	t.Run("Invalid SellerID", func(t *testing.T) {
		//Arrange
		service := NewMockService()
		server := createServer(&service)

		sellerId := "abcd"

		service.On("GetAllBySeller", sellerId).Return([]Product(nil), errors.New("The ID is not valid"))

		expectedHeaders := http.Header{
			"Content-Type": []string{"application/json; charset=utf-8"},
		}
		expectedBody := `{
			"error": "The ID is not valid"
		}`

		//Act
		request, respRecorder := createRequestTest(http.MethodGet, "/api/v1/products?seller_id=abcd", ``)
		server.ServeHTTP(respRecorder, request)

		//Assert
		assert.Equal(t, http.StatusInternalServerError, respRecorder.Code)
		assert.Equal(t, expectedHeaders, respRecorder.Header())
		assert.JSONEq(t, expectedBody, respRecorder.Body.String())
		service.AssertExpectations(t)
	})
}
