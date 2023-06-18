package repositories

import (
	"context"

	"tiu-85/gorest-iman/pkg/common/domain/entities"

	"tiu-85/gorest-iman/pkg/common/infra/adapters"
)

type TaskRepoFactory interface {
	Create(context.Context, adapters.DB) TaskRepo
}

type TaskRepo interface {
	FindById(uint32) (*entities.Task, error)
	Save(*entities.Task) error
}
