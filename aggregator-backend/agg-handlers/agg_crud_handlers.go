package agghandlers

import (
	"aggregator/database"
	"aggregator/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

// GetAllSubs godoc
// @Summary Get all subscriptions
// @Description Retrieves a list of all subscriptions from database
// @Tags subscriptions
// @Produce json
// @Success 200 {array} []models.UserSub "List of subscriptions"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /subscriptions [get]
func GetAllSubs(DB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var subs []models.UserSub

		if err := database.DB.Find(&subs).Error; err != nil {
			log.Printf("Failed to find subscriptions list: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get all subscriptions"})
			return
		}

		response := make([]models.SubscriptionResponse, len(subs))
		for i, sub := range subs {
			response[i] = models.ToSubscriptionResponse(sub)
		} //Для вывода даты в формате YYYY-MM

		c.JSON(http.StatusOK, response)
	}
}

// GetSubsById godoc
// @Summary Get subscription by ID
// @Description Retrieves a subscription by its ID.
// @Tags subscriptions
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} models.UserSub "Subscription details"
// @Failure 400 {object} map[string]string "Invalid ID format"
// @Failure 404 {object} map[string]string "Subscription not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /subscribe/{id} [get]
func GetSubsById(DB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var sub models.UserSub

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Printf("Invalid input id:%v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"Invalid id": err})
		}

		if err := database.DB.First(&sub, id).Error; err != nil {
			log.Printf("Failed to find subscription ID:%d, %v", sub.ID, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get subscription"})
			return
		}

		response := models.ToSubscriptionResponse(sub)

		c.JSON(http.StatusOK, response)
	}
}

// AddSubscription godoc
// @Summary Add a new subscription
// @Description Adds a new subscription to the database.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body models.SubscriptionInput true "Subscription data"
// @Success 201 {object} models.UserSub "Created subscription"
// @Failure 400 {object} map[string]string "Invalid input parameters or UUID format"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /subscribe [post]
func AddSubscription(DB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var sub models.SubscriptionInput

		if err := c.ShouldBindJSON(&sub); err != nil {
			log.Printf("Error to bind JSON:%v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"Invalid subscription input": err.Error()})
			return
		}

		userID, err := uuid.FromString(sub.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format"})
			return
		}

		startDate, err := time.Parse("2006-01", sub.StartDate)
		if err != nil {
			log.Printf("Invalid start_date format: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format, use YYYY-MM"})
			return
		}
		endDate, err := time.Parse("2006-01", sub.EndDate)
		if err != nil {
			log.Printf("Invalid end_date format: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format, use YYYY-MM"})
			return
		}

		// Проверка, что start_date <= end_date
		if startDate.After(endDate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "start_date must be before or equal to end_date"})
			return
		}

		newSub := models.UserSub{
			ServiceName: sub.ServiceName,
			UserID:      userID,
			Cost:        sub.Cost,
			StartDate:   startDate,
			EndDate:     endDate,
		}

		if err := database.DB.Create(&newSub).Error; err != nil {
			log.Printf("Failed to create subscription ID:%d, %v", newSub.ID, err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to create subscription"})
			return
		}

		response := models.ToSubscriptionResponse(newSub)

		c.JSON(http.StatusCreated, response)
	}
}

// UpdateSubscriptionById godoc
// @Summary Update a subscription
// @Description Updates a subscription by its ID. Creates a new one if not found.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Param subscription body models.SubscriptionInput true "Subscription data"
// @Success 200 {object} models.UserSub "Updated subscription"
// @Success 201 {object} models.UserSub "Created subscription"
// @Failure 400 {object} map[string]string "Invalid input parameters or UUID format"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /subscribe/{id} [put]
func UpdateSubscriptionById(DB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var sub models.SubscriptionInput
		var old_sub models.UserSub

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Printf("Invalid input id:%v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"Invalid id": err})
		}

		if err := c.ShouldBindJSON(&sub); err != nil {
			log.Printf("Error to bind JSON for subscription ID: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"Invalid subscription input": err.Error()})
			return
		}

		userID, err := uuid.FromString(sub.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format"})
			return
		}

		startDate, err := time.Parse("2006-01", sub.StartDate)
		if err != nil {
			log.Printf("Invalid start_date format: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format, use YYYY-MM"})
			return
		}
		endDate, err := time.Parse("2006-01", sub.EndDate)
		if err != nil {
			log.Printf("Invalid end_date format: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format, use YYYY-MM"})
			return
		}

		// Проверка, что start_date <= end_date
		if startDate.After(endDate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "start_date must be before or equal to end_date"})
			return
		}

		res := database.DB.First(&old_sub, id).Error
		if res != nil {
			var new_sub models.UserSub
			new_sub.ID = id
			new_sub.ServiceName = sub.ServiceName
			new_sub.UserID = userID
			new_sub.Cost = sub.Cost
			new_sub.StartDate = startDate
			new_sub.EndDate = endDate
			err := database.DB.Create(&new_sub).Error
			if err != nil {
				log.Printf("Failed to create subscription: %v", err)
				c.JSON(http.StatusBadRequest, gin.H{
					"Error to create subscription": err.Error()})
				return
			}
			var resp models.UserSub
			resp.ID = id
			resp.ServiceName = sub.ServiceName
			resp.UserID = userID
			resp.Cost = sub.Cost
			resp.StartDate = startDate
			resp.EndDate = endDate

			response := models.ToSubscriptionResponse(resp) //Так же форматируем дату
			c.JSON(http.StatusCreated, response)
			return
		}

		old_sub.ServiceName = sub.ServiceName
		old_sub.UserID = userID
		old_sub.Cost = sub.Cost
		old_sub.StartDate = startDate
		old_sub.EndDate = endDate

		if err := database.DB.Save(old_sub).Error; err != nil {
			log.Printf("Failed to update subscription by ID:%d, %v", old_sub.ID, err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to update subscription"})
			return
		}

		response := models.ToSubscriptionResponse(old_sub)

		c.JSON(http.StatusOK, response)
	}
}

// DeleteSubById godoc
// @Summary Delete a subscription
// @Description Deletes a subscription by its ID.
// @Tags subscriptions
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 204 "No content"
// @Failure 400 {object} map[string]string "Invalid ID format"
// @Failure 404 {object} map[string]string "Subscription not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /subscribe/{id} [delete]
func DeleteSubById(DB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var sub models.UserSub

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Printf("Invalid input id:%v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"Invalid id": err})
		}

		if err := database.DB.Delete(&sub, id).Error; err != nil {
			log.Printf("Failed to delete subscription by ID:%d, %v", sub.ID, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete subscription"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
