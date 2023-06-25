package entity

import (
	"gorm.io/gorm"
	"time"
)

type Movie struct {
	Id          uint64         `gorm:"column:id;type:bigint(20);AUTO_INCREMENT;primary_key"`
	Title       string         `gorm:"column:title;type:char(150);NOT NULL"`
	Description string         `gorm:"column:description;type:text;NOT NULL"`
	Duration    string         `gorm:"column:duration;type:varchar(250);NOT NULL"`
	Artist      string         `gorm:"column:artist;type:char(50);NOT NULL"`
	Genre       string         `gorm:"column:genre;type:char(50);NOT NULL"`
	VideoUrl    string         `gorm:"column:video_url;type:text;NOT NULL"`
	CreatedAt   time.Time      `gorm:"column:created_at;type:timestamp;NOT NULL"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;type:timestamp;NOT NULL"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp"`
}

func (e *Movie) TableName() string {
	return "movies"
}
