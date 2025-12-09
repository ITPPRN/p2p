package models

import (
	"time"

	"gorm.io/gorm"
)

type UserEntity struct {
	ID        string         `gorm:"primaryKey;type:char(36)" json:"id"` // ใช้ char(36) สำหรับ UUID
	Username  string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
	Email     string         `gorm:"type:varchar(150);uniqueIndex;not null" json:"email"`
	FirstName string         `gorm:"type:varchar(100)" json:"first_name"`
	LastName  string         `gorm:"type:varchar(100)" json:"last_name"`
	Role      string         `gorm:"type:varchar(50)" json:"role"`
	Password  string         `gorm:"type:varchar(255);not null" json:"password"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}


