package orchestrator

import (
	"github.com/VaDKustiK/yandex-golang-course/calculator_service/pkg/common"
	"gorm.io/gorm"
)

var db *gorm.DB

func SetDB(d *gorm.DB) {
	db = d
}

func GetDB() *gorm.DB {
	return db
}

func getExpressionByID(id uint) (*common.Expression, error) {
	var expr common.Expression
	err := db.Preload("Tasks").
		Where("id = ?", id).
		First(&expr).Error
	return &expr, err
}
