package models

import (
	"time"
)

type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;default:1" json:"user_id"`
	Name      string    `gorm:"not null" json:"name"`
	Color     string    `json:"color"`
	SortOrder int       `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Activity struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;default:1" json:"user_id"`
	CategoryID uint      `gorm:"not null" json:"category_id"`
	Category   Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Name       string    `gorm:"not null" json:"name"`
	Pinned     int       `gorm:"not null;default:0" json:"pinned"` // 0 or 1
	Color      string    `json:"color"`
	SortOrder  int       `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type TimeEntry struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	UserID     uint       `gorm:"not null;default:1" json:"user_id"`
	CategoryID uint       `gorm:"not null" json:"category_id"`
	Category   Category   `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	ActivityID uint       `gorm:"not null" json:"activity_id"`
	Activity   Activity   `gorm:"foreignKey:ActivityID" json:"activity,omitempty"`
	StartTime  time.Time  `gorm:"not null;index:idx_time_entries_user_start" json:"start_time"`
	EndTime    *time.Time `gorm:"index:idx_time_entries_user_end" json:"end_time"` // Nullable for running entries
	Note       string     `json:"note"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
