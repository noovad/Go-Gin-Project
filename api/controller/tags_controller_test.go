package controller_test

import (
	"bytes"
	"go-gin-project/api/controller"
	"go-gin-project/data"
	"go-gin-project/helper"
	"net/http"
	"net/http/httptest"
	"testing"

	"errors"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTagsService struct {
	mock.Mock
}

func (m *MockTagsService) Create(request data.TagRequest) error {
	args := m.Called(request)
	return args.Error(0)
}

func (m *MockTagsService) FindAll() ([]data.TagResponse, error) {
	args := m.Called()
	return args.Get(0).([]data.TagResponse), args.Error(1)
}

func (m *MockTagsService) FindById(tagId string) (data.TagResponse, error) {
	args := m.Called(tagId)
	return args.Get(0).(data.TagResponse), args.Error(1)
}

func (m *MockTagsService) Update(tagId string, request data.TagRequest) error {
	args := m.Called(tagId, request)
	return args.Error(0)
}

func (m *MockTagsService) Delete(tagId int) error {
	args := m.Called(tagId)
	return args.Error(0)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	return router
}

func setupTest() (*MockTagsService, *controller.TagsController, *gin.Engine) {
	mockService := new(MockTagsService)
	controller := controller.NewTagsController(mockService)
	router := setupRouter()
	return mockService, controller, router
}

// Create tests
func TestCreateTag(t *testing.T) {
	t.Run("success - should create tag", func(t *testing.T) {
		mockService, controller, router := setupTest()
		router.POST("/tags", controller.Create)

		mockService.On("Create", mock.Anything).Return(nil)

		requestBody := `{"name": "New Tag"}`
		req, _ := http.NewRequest("POST", "/tags", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), `"status":"Successfully created"`)
		mockService.AssertExpectations(t)
	})

	t.Run("error - should return internal server error", func(t *testing.T) {
		mockService, controller, router := setupTest()
		router.POST("/tags", controller.Create)

		mockService.On("Create", mock.Anything).Return(errors.New("unexpected error"))

		requestBody := `{"name": "New Tag"}`
		req, _ := http.NewRequest("POST", "/tags", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("error - should handle invalid JSON", func(t *testing.T) {
		_, controller, router := setupTest()
		router.POST("/tags", controller.Create)

		requestBody := `{"name":`
		req, _ := http.NewRequest("POST", "/tags", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("error - should handle validation error", func(t *testing.T) {
		mockService, controller, router := setupTest()
		router.POST("/tags", controller.Create)

		mockService.On("Create", mock.Anything).Return(helper.ErrFailedValidation)

		requestBody := `{}`
		req, _ := http.NewRequest("POST", "/tags", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})
}

// FindAll tests
func TestFindAllTags(t *testing.T) {
	t.Run("success - should return all tags", func(t *testing.T) {
		mockService, controller, router := setupTest()
		router.GET("/tags", controller.FindAll)

		expectedTags := []data.TagResponse{
			{Id: 1, Name: "Tag1"},
			{Id: 2, Name: "Tag2"},
		}
		mockService.On("FindAll").Return(expectedTags, nil)

		req, _ := http.NewRequest("GET", "/tags", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"name":"Tag1"`)
		mockService.AssertExpectations(t)
	})

	t.Run("error - should handle server error", func(t *testing.T) {
		mockService, controller, router := setupTest()
		router.GET("/tags", controller.FindAll)

		mockService.On("FindAll").Return([]data.TagResponse{}, errors.New("unexpected error"))

		req, _ := http.NewRequest("GET", "/tags", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

// FindById tests
func TestFindById(t *testing.T) {
	t.Run("success - should return tag by ID", func(t *testing.T) {
		mockService, controller, router := setupTest()
		router.GET("/tags/:tagId", controller.FindById)

		expectedTag := data.TagResponse{Id: 1, Name: "Tag1"}
		mockService.On("FindById", "1").Return(expectedTag, nil)

		req, _ := http.NewRequest("GET", "/tags/1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"name":"Tag1"`)
		mockService.AssertExpectations(t)
	})

	t.Run("error - should return not found", func(t *testing.T) {
		mockService, controller, router := setupTest()
		router.GET("/tags/:tagId", controller.FindById)

		mockService.On("FindById", "999").Return(data.TagResponse{}, helper.ErrNotFound)

		req, _ := http.NewRequest("GET", "/tags/999", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("error - should return internal server error", func(t *testing.T) {
		mockService, controller, router := setupTest()
		router.GET("/tags/:tagId", controller.FindById)

		mockService.On("FindById", "1").Return(data.TagResponse{}, errors.New("unexpected error"))

		req, _ := http.NewRequest("GET", "/tags/1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), `"status":"Internal Server Error"`)
		mockService.AssertExpectations(t)
	})
}

// Update tests
func TestUpdateTag(t *testing.T) {
	t.Run("success - should update tag", func(t *testing.T) {
		mockService, controller, router := setupTest()
		router.PUT("/tags/:tagId", controller.Update)

		requestBody := `{"name": "Updated Tag"}`
		req, _ := http.NewRequest("PUT", "/tags/1", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockService.On("Update", "1", mock.Anything).Return(nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"status":"Successfully updated"`)
		mockService.AssertExpectations(t)
	})

	t.Run("error - should handle validation error", func(t *testing.T) {
		mockService, controller, router := setupTest()
		router.PUT("/tags/:tagId", controller.Update)

		requestBody := `{"name": ""}`
		req, _ := http.NewRequest("PUT", "/tags/1", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockService.On("Update", "1", mock.Anything).Return(helper.ErrFailedValidation)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"status":"Bad Request"`)
		mockService.AssertExpectations(t)
	})

	t.Run("error - should return not found", func(t *testing.T) {
		mockService, controller, router := setupTest()
		router.PUT("/tags/:tagId", controller.Update)

		requestBody := `{"name": "Updated Tag"}`
		req, _ := http.NewRequest("PUT", "/tags/999", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockService.On("Update", "999", mock.Anything).Return(helper.ErrNotFound)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), `"status":"Not Found"`)
		mockService.AssertExpectations(t)
	})

	t.Run("error - should handle invalid JSON", func(t *testing.T) {
		mockService, controller, router := setupTest()
		router.PUT("/tags/:tagId", controller.Update)

		requestBody := `{"name":`
		req, _ := http.NewRequest("PUT", "/tags/1", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), `"status":"Internal Server Error"`)
		mockService.AssertExpectations(t)
	})

	t.Run("error - should return internal server error", func(t *testing.T) {
		mockService, controller, router := setupTest()
		router.PUT("/tags/:tagId", controller.Update)

		requestBody := `{"name": "Updated Tag"}`
		req, _ := http.NewRequest("PUT", "/tags/1", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockService.On("Update", "1", mock.Anything).Return(errors.New("unexpected error"))

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), `"status":"Internal Server Error"`)
		mockService.AssertExpectations(t)
	})
}

// Delete tests
func TestDeleteTag(t *testing.T) {
	t.Run("success - should delete tag", func(t *testing.T) {
		mockService, controller, router := setupTest()
		router.DELETE("/tags/:tagId", controller.Delete)

		mockService.On("Delete", 1).Return(nil)

		req, _ := http.NewRequest("DELETE", "/tags/1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"status":"Successfully deleted"`)
		mockService.AssertExpectations(t)
	})

	t.Run("error - should return not found", func(t *testing.T) {
		mockService, controller, router := setupTest()
		router.DELETE("/tags/:tagId", controller.Delete)

		mockService.On("Delete", 999).Return(helper.ErrNotFound)

		req, _ := http.NewRequest("DELETE", "/tags/999", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("error - should handle invalid ID", func(t *testing.T) {
		mockService, controller, router := setupTest()
		router.DELETE("/tags/:tagId", controller.Delete)

		req, _ := http.NewRequest("DELETE", "/tags/invalid-id", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"status":"Bad Request"`)
		mockService.AssertExpectations(t)
	})

	t.Run("error - should return internal server error", func(t *testing.T) {
		mockService, controller, router := setupTest()
		router.DELETE("/tags/:tagId", controller.Delete)

		mockService.On("Delete", 1).Return(errors.New("unexpected error"))

		req, _ := http.NewRequest("DELETE", "/tags/1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), `"status":"Internal Server Error"`)
		mockService.AssertExpectations(t)
	})
}
