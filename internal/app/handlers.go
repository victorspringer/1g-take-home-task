package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *handler) healthCheck() func(c *gin.Context) {
	return func(c *gin.Context) {
		h.logger.Debug("health check", zap.String("requestUrl", c.Request.URL.Path))
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	}
}

func (h *handler) listAllDevices() func(c *gin.Context) {
	return func(c *gin.Context) {
		h.logger.Debug("list all devices", zap.String("requestUrl", c.Request.URL.Path))

		id := c.Param("id")

		c.JSON(http.StatusOK, gin.H{
			"id": id,
		})
	}
}

func (h *handler) getDeviceByID() func(c *gin.Context) {
	return func(c *gin.Context) {
		h.logger.Debug("get device by id", zap.String("requestUrl", c.Request.URL.Path))

		id := c.Param("id")

		c.JSON(http.StatusOK, gin.H{
			"id": id,
		})
	}
}

func (h *handler) searchDevices() func(c *gin.Context) {
	return func(c *gin.Context) {
		h.logger.Debug("search device", zap.String("requestUrl", c.Request.URL.Path+"?"+c.Request.URL.Query().Encode()))

		brand := c.Query("brand")

		c.JSON(http.StatusOK, gin.H{
			"brand": brand,
		})
	}
}

func (h *handler) addDevice() func(c *gin.Context) {
	return func(c *gin.Context) {
		h.logger.Debug("add device")

		c.JSON(http.StatusOK, gin.H{
			"action": "add device",
		})
	}
}

func (h *handler) updateDevice() func(c *gin.Context) {
	return func(c *gin.Context) {
		h.logger.Debug("update device", zap.String("requestUrl", c.Request.URL.Path))

		id := c.Param("id")

		c.JSON(http.StatusOK, gin.H{
			"id": id,
		})
	}
}

func (h *handler) deleteDevice() func(c *gin.Context) {
	return func(c *gin.Context) {
		h.logger.Debug("delete device", zap.String("requestUrl", c.Request.URL.Path))

		id := c.Param("id")

		c.JSON(http.StatusOK, gin.H{
			"id": id,
		})
	}
}
