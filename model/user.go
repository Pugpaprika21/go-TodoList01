package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(50);not null"`
	Password string `gorm:"type:varchar(50);not null"`
	Status   int
	Todos    []Todo // เพิ่มสัมพันธ์ One-to-Many ระหว่าง User และ Todo
}
