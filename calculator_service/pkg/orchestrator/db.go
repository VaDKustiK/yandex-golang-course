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
	if err := db.Preload("Tasks").First(&expr, id).Error; err != nil {
		return nil, err
	}
	return &expr, nil
}
