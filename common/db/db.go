package db

import "github.com/jinzhu/gorm"

var (
	db *gorm.DB
)

func InitDb() {

}

func Db() *gorm.DB {
	return db
}
