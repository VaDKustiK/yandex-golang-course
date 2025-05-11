package orchestrator

import (
	"context"
	"errors"

	"github.com/VaDKustiK/yandex-golang-course/calculator_service/pkg/common"
	"gorm.io/gorm"
)

var ErrNoTask = errors.New("no pending task")

func FetchNextPendingTask(ctx context.Context) (*common.Task, error) {
	t := new(common.Task)
	if err := db.Where("status = ?", "pending").First(t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoTask
		}
		return nil, err
	}
	if err := db.Model(t).Update("status", "in_progress").Error; err != nil {
		return nil, err
	}
	return t, nil
}

func StoreTaskResult(ctx context.Context, id uint, result float64) error {
	t := new(common.Task)
	if err := db.First(t, id).Error; err != nil {
		return err
	}
	return db.Model(t).Updates(map[string]interface{}{
		"status": "completed",
		"result": result,
	}).Error
}
