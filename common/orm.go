package common

import (
	"gorm.io/gorm"
)

func SoftDelete(db *gorm.DB) *gorm.DB {
	return db.Where("deleted_at IS NULL")
}