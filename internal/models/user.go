package models

import (
	"time"
)

type User struct {
	ID                  uint       `json:"id" gorm:"primaryKey"`
	Email               string     `json:"email" gorm:"uniqueIndex;not null"`
	Password            string     `json:"-" gorm:"not null"`
	IsActive            bool       `json:"is_active" gorm:"default:false"`
	EmailVerified       bool       `json:"email_verified" gorm:"default:false"`
	EmailVerifyToken    string     `json:"-" gorm:"index"`
	EmailVerifyExpires  *time.Time `json:"-"`
	DailyLimit          int        `json:"daily_limit" gorm:"default:1"`
	MonthlyLimit        int        `json:"monthly_limit" gorm:"default:30"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type Request struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      *uint     `json:"user_id" gorm:"index"`
	URL         string    `json:"url" gorm:"not null"`
	Email       string    `json:"email" gorm:"not null"`
	Format      string    `json:"format" gorm:"default:html"`
	Status      string    `json:"status" gorm:"default:pending"`
	FilePath    string    `json:"file_path"`
	ErrorMsg    string    `json:"error_msg"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
	
	User        *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type EmailConfig struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Host     string `json:"host" gorm:"not null"`
	Port     int    `json:"port" gorm:"not null"`
	Username string `json:"username" gorm:"not null"`
	Password string `json:"-" gorm:"not null"`
	FromName string `json:"from_name"`
	IsActive bool   `json:"is_active" gorm:"default:true"`
}

type EmailVerification struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Email      string    `json:"email" gorm:"not null;index"`
	IPAddress  string    `json:"ip_address" gorm:"not null"`
	SentAt     time.Time `json:"sent_at"`
	CreatedAt  time.Time `json:"created_at"`
}