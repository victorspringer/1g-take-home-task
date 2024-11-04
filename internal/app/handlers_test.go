package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/victorspringer/1g-take-home-task/internal/pkg/device"
	"go.uber.org/zap"
)

func setupRouter(repo device.Repository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	h := &handler{
		logger:           zap.NewNop(),
		deviceRepository: repo,
	}

	router.GET("/devices", h.listAllDevices)
	router.GET("/devices/:id", h.getDeviceByID)
	router.GET("/devices/search", h.searchDevices)
	router.POST("/devices", h.addDevice)
	router.PATCH("/devices/:id", h.updateDevice)
	router.DELETE("/devices/:id", h.deleteDevice)

	return router
}

func TestListAllDevices_Success(t *testing.T) {
	repo := &device.MockRepository{
		Devices: []device.Device{
			{ID: "1", Name: "Device1", Brand: "BrandA"},
		},
	}
	router := setupRouter(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/devices", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	expectedResponse, _ := json.Marshal(gin.H{"devices": repo.Devices})
	assert.JSONEq(t, string(expectedResponse), w.Body.String())
}

func TestListAllDevices_NoDevices(t *testing.T) {
	repo := &device.MockRepository{}
	router := setupRouter(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/devices", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, `{"status": "Not Found"}`, w.Body.String())
}

func TestListAllDevices_Error(t *testing.T) {
	repo := &device.MockRepository{
		Err: errors.New("internal error"),
	}
	router := setupRouter(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/devices", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"error": "internal error"}`, w.Body.String())
}

func TestGetDeviceByID_Success(t *testing.T) {
	repo := &device.MockRepository{
		Devices: []device.Device{
			{ID: "1", Name: "Device1", Brand: "BrandA"},
		},
	}
	router := setupRouter(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/devices/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	expectedResponse, _ := json.Marshal(gin.H{"device": repo.Devices[0]})
	assert.JSONEq(t, string(expectedResponse), w.Body.String())
}

func TestGetDeviceByID_NotFound(t *testing.T) {
	repo := &device.MockRepository{
		Devices: []device.Device{
			{ID: "1", Name: "Device1", Brand: "BrandA"},
		},
	}
	router := setupRouter(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/devices/2", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, `{"status": "Not Found"}`, w.Body.String())
}

func TestGetDeviceByID_Error(t *testing.T) {
	repo := &device.MockRepository{
		Err: errors.New("internal error"),
	}
	router := setupRouter(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/devices/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"error": "internal error"}`, w.Body.String())
}

func TestSearchDevices_Success(t *testing.T) {
	repo := &device.MockRepository{
		Devices: []device.Device{
			{ID: "1", Name: "Device1", Brand: "BrandA"},
		},
	}
	router := setupRouter(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/devices/search?brand=BrandA", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	expectedResponse, _ := json.Marshal(gin.H{"devices": repo.Devices})
	assert.JSONEq(t, string(expectedResponse), w.Body.String())
}

func TestSearchDevices_NoResults(t *testing.T) {
	repo := &device.MockRepository{}
	router := setupRouter(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/devices/search?brand=NonExistent", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, `{"status": "Not Found"}`, w.Body.String())
}

func TestSearchDevices_Error(t *testing.T) {
	repo := &device.MockRepository{
		Err: errors.New("internal error"),
	}
	router := setupRouter(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/devices/search?brand=BrandA", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"error": "internal error"}`, w.Body.String())
}

func TestAddDevice_Success(t *testing.T) {
	repo := &device.MockRepository{}
	router := setupRouter(repo)

	newDevice := device.Device{Name: "Device1", Brand: "BrandA"}
	jsonDevice, _ := json.Marshal(newDevice)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/devices", bytes.NewBuffer(jsonDevice))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, `{"status": "Created"}`, w.Body.String())
}

func TestAddDevice_BadRequest(t *testing.T) {
	repo := &device.MockRepository{}
	router := setupRouter(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/devices", bytes.NewBuffer([]byte("{")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error": "unexpected EOF"}`, w.Body.String())
}

func TestAddDevice_Error(t *testing.T) {
	repo := &device.MockRepository{
		Err: errors.New("internal error"),
	}
	router := setupRouter(repo)

	newDevice := device.Device{Name: "Device1", Brand: "BrandA"}
	jsonDevice, _ := json.Marshal(newDevice)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/devices", bytes.NewBuffer(jsonDevice))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"error": "internal error"}`, w.Body.String())
}

func TestUpdateDevice_Success(t *testing.T) {
	repo := &device.MockRepository{
		Devices: []device.Device{
			{ID: "1", Name: "Device1", Brand: "BrandA"},
		},
	}
	router := setupRouter(repo)

	updatedDevice := device.Device{Name: "UpdatedDevice", Brand: "BrandB"}
	jsonDevice, _ := json.Marshal(updatedDevice)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/devices/1", bytes.NewBuffer(jsonDevice))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status": "OK"}`, w.Body.String())
}

func TestUpdateDevice_Error(t *testing.T) {
	repo := &device.MockRepository{
		Err: errors.New("internal error"),
	}
	router := setupRouter(repo)

	updatedDevice := device.Device{Name: "UpdatedDevice", Brand: "BrandB"}
	jsonDevice, _ := json.Marshal(updatedDevice)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/devices/1", bytes.NewBuffer(jsonDevice))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"error": "internal error"}`, w.Body.String())
}

func TestDeleteDevice_Success(t *testing.T) {
	repo := &device.MockRepository{
		Devices: []device.Device{
			{ID: "1", Name: "Device1", Brand: "BrandA"},
		},
	}
	router := setupRouter(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/devices/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteDevice_Error(t *testing.T) {
	repo := &device.MockRepository{
		Err: errors.New("internal error"),
	}
	router := setupRouter(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/devices/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"error": "internal error"}`, w.Body.String())
}
