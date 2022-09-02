package entity

import "time"

const UserTableName = "user_manager_User"

type User struct {
	Id          int    `json:"id" form:"id"`
	UserID      string `json:"user_id" form:"user_id"`
	Version     string `json:"version" form:"version"`
	UserName    string `json:"user_name" form:"user_name"`
	Description string `json:"description" form:"description"`
	Template    string `json:"template" form:"-"`
	ConfigName  string `json:"-" form:"config_name"`
}

type UserListQuery struct {
	UserName string `form:"user_name"`
	Sort     string `form:"sort" validate:"omitempty,oneof=id User_name"`
	CommonListQuery
}

type Filter struct {
	Username    string    `json:"username" binding:"required"`
	Category    string    `json:"category" binding:"required"`
	StartTime   int       `json:"starttime" binding:"gte=0"`
	EndTime     int       `json:"endtime" binding:"gte=0"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}
