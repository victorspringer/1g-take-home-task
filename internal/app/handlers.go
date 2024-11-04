package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victorspringer/1g-take-home-task/internal/pkg/device"
	"go.uber.org/zap"
)

type handler struct {
	logger           *zap.Logger
	deviceRepository device.Repository
}

func (h *handler) healthCheck(c *gin.Context) {
	h.logger.Debug("health check", zap.String("requestUrl", c.Request.URL.Path))
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusText(http.StatusOK),
	})
}

func (h *handler) listAllDevices(c *gin.Context) {
	h.logger.Debug("list all devices", zap.String("requestUrl", c.Request.URL.Path))

	devices, err := h.deviceRepository.List()
	if err != nil {
		h.logger.Error("error listing all devices", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(devices) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusText(http.StatusNotFound),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"devices": devices,
	})
}

func (h *handler) getDeviceByID(c *gin.Context) {
	h.logger.Debug("get device by id", zap.String("requestUrl", c.Request.URL.Path))

	id := c.Param("id")
	device, err := h.deviceRepository.FindByID(id)
	if err != nil {
		h.logger.Error("error getting device by id", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if device == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusText(http.StatusNotFound),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"device": device,
	})
}

func (h *handler) searchDevices(c *gin.Context) {
	h.logger.Debug("search device", zap.String("requestUrl", c.Request.URL.Path+"?"+c.Request.URL.Query().Encode()))

	brand := c.Query("brand")

	devices, err := h.deviceRepository.FindByBrand(brand)
	if err != nil {
		h.logger.Error("error searching devices by brand", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(devices) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusText(http.StatusNotFound),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"devices": devices,
	})
}

func (h *handler) addDevice(c *gin.Context) {
	h.logger.Debug("add device")

	var dvc device.Device

	if err := c.ShouldBindJSON(&dvc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if dvc.ID != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "device id is not a valid field",
		})
		return
	}

	if err := h.deviceRepository.Store(&dvc); err != nil {
		h.logger.Error("error adding device", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": http.StatusText(http.StatusCreated),
	})
}

func (h *handler) updateDevice(c *gin.Context) {
	h.logger.Debug("update device", zap.String("requestUrl", c.Request.URL.Path))

	id := c.Param("id")

	var dvc device.Device

	if err := c.ShouldBindJSON(&dvc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	dvc.ID = id

	if err := h.deviceRepository.Update(&dvc); err != nil {
		h.logger.Error("error updating device", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusText(http.StatusOK),
	})
}

func (h *handler) deleteDevice(c *gin.Context) {
	h.logger.Debug("delete device", zap.String("requestUrl", c.Request.URL.Path))

	id := c.Param("id")

	if err := h.deviceRepository.Remove(id); err != nil {
		h.logger.Error("error deleting device", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusText(http.StatusOK),
	})
}
