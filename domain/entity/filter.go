package entity

import "time"

type Filter struct {
	Username    string    `json:"username" binding:"required"`
	Category    string    `json:"category" binding:"required"`
	StartTime   int       `json:"starttime" binding:"gte=0"`
	EndTime     int       `json:"endtime" binding:"gte=0"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}
