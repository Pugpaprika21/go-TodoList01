package model

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	Title       string `gorm:"type:varchar(100);not null"`
	Description string `gorm:"type:varchar(255);not null"`
	UserID      uint   // คีย์ UserID เพื่อเชื่อมโยงกับผู้ใช้ที่สร้าง Todo
}
