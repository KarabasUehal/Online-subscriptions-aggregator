package models

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
)

type UserSub struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ServiceName string    `json:"service_name" gorm:"type:varchar(255);not null"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;index"`
	Cost        int       `json:"cost" gorm:"type:integer;not null"`
	StartDate   time.Time `json:"start_date" gorm:"type:date;not null"`
	EndDate     time.Time `json:"end_date" gorm:"type:date;not null"`
}

func (s *UserSub) ValidateDates() error {
	const layout = "2006-01"
	if _, err := time.Parse(layout, s.StartDate.String()); err != nil {
		return fmt.Errorf("invalid start_date format: %v, format must be YYYY-MM", err)
	}
	if _, err := time.Parse(layout, s.EndDate.String()); err != nil {
		return fmt.Errorf("invalid end_date format: %v, format must be YYYY-MM", err)
	}
	return nil
}

type SubscriptionInput struct {
	UserID      string `json:"user_id" binding:"required,uuid"`
	ServiceName string `json:"service_name" binding:"required"`
	Cost        int    `json:"cost" binding:"required,gte=0"`
	StartDate   string `json:"start_date" binding:"required"`
	EndDate     string `json:"end_date" binding:"required"`
}

type SubscriptionResponse struct {
	ID          int       `json:"id"`
	ServiceName string    `json:"service_name"`
	Cost        int       `json:"cost"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     string    `json:"end_date"`
}

func ToSubscriptionResponse(sub UserSub) SubscriptionResponse {
	return SubscriptionResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Cost:        sub.Cost,
		UserID:      sub.UserID,
		StartDate:   sub.StartDate.Format("2006-01"), // Формат YYYY-MM, как просили в тз
		EndDate:     sub.EndDate.Format("2006-01"),
	}
}
