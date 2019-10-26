package database

import (
	"github.com/jinzhu/gorm"
)

var (
	DBCon *gorm.DB
	Error error
)
