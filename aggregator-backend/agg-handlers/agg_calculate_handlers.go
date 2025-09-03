package agghandlers

import (
	"aggregator/database"
	"aggregator/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

// GetTotalCost godoc
// @Summary Calculate total cost of subscriptions
// @Description Calculates the total cost of subscriptions for a given period, optionally filtered by user ID and service name.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id query string false "User ID (UUID format)"
// @Param service_name query string false "Service name"
// @Param start query string true "Start date of the period (YYYY-MM)"
// @Param end query string true "End date of the period (YYYY-MM)"
// @Success 200 {object} map[string]float64 "Total cost"
// @Failure 400 {object} map[string]string "Invalid input parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /total-cost [get]
func CalculateSubsPerPeriod(DB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Query("user_id")
		serviceName := c.Query("service_name")
		startDateStr := c.Query("start_date")
		endDateStr := c.Query("end_date")

		if startDateStr == "" || endDateStr == "" {
			log.Printf("Invalid input date")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "start_date and end_date are required",
			})
			return
		}

		startDate, err := time.Parse("2006-01", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format (YYYY-MM-DD)"})
			return
		}
		endDate, err := time.Parse("2006-01", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format (YYYY-MM-DD)"})
			return
		}
		if startDate.After(endDate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "start_date must be before end_date"})
			return
		}

		var userID uuid.UUID
		if userIDStr != "" {
			userID, err = uuid.FromString(userIDStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format (UUID)"})
				return
			}
		}

		query := database.DB.Model(&models.UserSub{})

		// Здесь фильтры
		if userIDStr != "" {
			query = query.Where("user_id = ?", userID)
		}
		if serviceName != "" {
			query = query.Where("service_name = ?", serviceName)
		}

		query = query.Where("start_date <= ? AND end_date >= ?", endDate, startDate)

		var subs []models.UserSub
		if err := query.Find(&subs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query subs"})
			return
		}

		// Подсчитывание общей стоимости
		totalCost := 0
		for _, sub := range subs {
			// Подсчёт периода
			overlapStart := maxTime(startDate, sub.StartDate)
			overlapEnd := minTime(endDate, sub.EndDate)

			// Пропуск , если период не совпадает
			if overlapStart.After(overlapEnd) {
				continue
			}

			// Подсчёт длительности в днях (входящие данные выравнены по полуночи)
			overlapDays := overlapEnd.Sub(overlapStart).Hours() / 24
			totalDays := sub.EndDate.Sub(sub.StartDate).Hours() / 24

			if totalDays <= 0 {
				continue // Невалидная длительность подписки
			}

			// Стоимость за все дни
			proportional := (int(overlapDays) / int(totalDays)) * sub.Cost
			totalCost += proportional
		}

		c.JSON(http.StatusOK, gin.H{"total_cost": totalCost})

	}
}

func maxTime(t1, t2 time.Time) time.Time {
	if t1.After(t2) {
		return t1
	}
	return t2
}

func minTime(t1, t2 time.Time) time.Time {
	if t1.Before(t2) {
		return t1
	}
	return t2
}
