package services

import (
	"context"

	"tiu-85/gorest-iman/pkg/common/domain/entities"
	"tiu-85/gorest-iman/pkg/common/domain/repositories"
	"tiu-85/gorest-iman/pkg/common/gen/pbv1"
	"tiu-85/gorest-iman/pkg/common/infra/adapters"
	"tiu-85/gorest-iman/pkg/common/infra/values"
	pdp "tiu-85/gorest-iman/pkg/post-fetch/domains/processors"
)

type PostFetchService struct {
	ctx              context.Context
	logger           adapters.Logger
	db               adapters.DB
	postApiProcessor pdp.PostApiProcessor
	taskRepoFactory  repositories.TaskRepoFactory
}

func NewPostFetchService(
	ctx context.Context,
	logger adapters.Logger,
	pool adapters.Pool,
	postApiProcessor pdp.PostApiProcessor,
	taskRepoFactory repositories.TaskRepoFactory,
) (*PostFetchService, error) {
	db, err := pool.Get(adapters.PostKind)
	if err != nil {
		return nil, err
	}

	return &PostFetchService{
		ctx:              ctx,
		logger:           logger,
		db:               db,
		postApiProcessor: postApiProcessor,
		taskRepoFactory:  taskRepoFactory,
	}, nil
}

func (s *PostFetchService) GetTask(ctx context.Context, req *pbv1.GetTaskRequest) (*pbv1.GetTaskResponse, error) {
	logger := s.logger.WithCtx(ctx, "method", "get_task")
	taskRepo := s.taskRepoFactory.Create(ctx, s.db)

	task, err := taskRepo.FindById(req.Id)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &pbv1.GetTaskResponse{
		Task: task.ToProto(),
	}, nil
}

func (s *PostFetchService) CreateTask(ctx context.Context, req *pbv1.CreateTaskRequest) (*pbv1.CreateTaskResponse, error) {
	logger := s.logger.WithCtx(ctx, "method", "create_task")

	if req.Paginator == nil {
		req.Paginator = &pbv1.Paginator{
			Offset: 0,
			Limit:  10,
		}
	}

	if req.Paginator.Offset < 0 || req.Paginator.Limit < 0 {
		logger.Errorf("invalid offset or limit")
		return nil, values.ErrInvalidParams
	}

	taskRepo := s.taskRepoFactory.Create(ctx, s.db)
	task := &entities.Task{
		Total:      req.Paginator.Offset + req.Paginator.Limit,
		PageOffset: req.Paginator.Offset,
		PageLimit:  req.Paginator.Limit,
	}

	err := taskRepo.Save(task)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	s.postApiProcessor.Process(task)

	return &pbv1.CreateTaskResponse{
		Task: task.ToProto(),
	}, nil
}
