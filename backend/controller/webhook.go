package controllers

import (
	"log"
	"net/http"

	"bus-tracking-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WebhookController struct {
	DB *gorm.DB
}

type SubscribeWebhookRequest struct {
	URL string `json:"url" binding:"required,url"`
}

func (wc *WebhookController) SubscribeWebhook(c *gin.Context) {
	var req SubscribeWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	webhook := models.Webhook{URL: req.URL}
	if err := wc.DB.Create(&webhook).Error; err != nil {
		if isUniqueConstraintError(err) {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "Webhook URL already subscribed",
			})
			return
		}
		log.Printf("DB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to subscribe webhook",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Webhook subscribed successfully",
		"data":    webhook,
	})
}

func isUniqueConstraintError(err error) bool {
	return false
}
