package controllers

import (
	"log"
	"net/http"
	"strconv"

	"bus-tracking-backend/services"

	"github.com/gin-gonic/gin"
)

type ETAController struct{}

func (ec *ETAController) GetETA(c *gin.Context) {
	busIDStr := c.Query("bus_id")
	stopIDStr := c.Query("stop_id")

	if busIDStr == "" || stopIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "bus_id and stop_id are required",
		})
		return
	}

	busID, err := strconv.ParseUint(busIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid bus_id",
		})
		return
	}

	stopID, err := strconv.ParseUint(stopIDStr, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid stop_id",
		})
		return
	}

	eta, err := services.CalculateETA(uint(busID), uint(stopID))
	if err != nil {
		log.Printf("ETA error for bus %d at stop %d: %v", busID, stopID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to calculate ETA",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "ETA calculated successfully",
		"data": map[string]interface{}{
			"bus_id":      busID,
			"stop_id":     stopID,
			"eta_minutes": eta,
		},
	})
}
