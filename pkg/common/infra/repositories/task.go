package repositories

import (
	"context"
	"errors"

	"github.com/go-pg/pg/v9"

	"tiu-85/gorest-iman/pkg/common/domain/entities"

	"tiu-85/gorest-iman/pkg/common/domain/repositories"
	"tiu-85/gorest-iman/pkg/common/infra/adapters"
	"tiu-85/gorest-iman/pkg/common/infra/values"
)

func NewTaskRepoFactory(
	ctx context.Context,
	logger adapters.Logger,
) repositories.TaskRepoFactory {
	return &taskRepoFactory{
		ctx:    ctx,
		logger: logger,
	}
}

type taskRepoFactory struct {
	ctx    context.Context
	logger adapters.Logger
}

func (f *taskRepoFactory) Create(ctx context.Context, db adapters.DB) repositories.TaskRepo {
	return &taskRepo{
		DS:     adapters.NewDS(db.WithCtx(ctx)),
		logger: f.logger.WithCtx(ctx, "repo", "task"),
	}
}

type taskRepo struct {
	*adapters.DS
	logger adapters.Logger
}

func (r *taskRepo) FindById(id uint32) (*entities.Task, error) {
	logger := r.logger.With("method", "find_by_id")
	model := new(entities.Task)

	err := r.Model(model).Where("task.id = ?", id).Select()
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			logger.Warn(err)
			return nil, values.ErrRowNotFound
		}

		logger.Error(err)
		return nil, err
	}

	return model, err
}

func (r *taskRepo) Save(model *entities.Task) error {
	logger := r.logger.With("method", "save")

	var err error
	if model.Id == 0 {
		err = r.Insert(model)
	} else {
		err = r.Update(model)
	}

	if err != nil {
		logger.Error(err)
		return values.ErrCantSave
	}

	return nil
}
