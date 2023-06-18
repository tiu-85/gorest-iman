package processors

import (
	"context"
	"sync"

	"tiu-85/gorest-iman/pkg/common/domain/entities"
	"tiu-85/gorest-iman/pkg/common/domain/repositories"
	"tiu-85/gorest-iman/pkg/common/infra/adapters"
	"tiu-85/gorest-iman/pkg/post-fetch/domains/clients"
	pdc "tiu-85/gorest-iman/pkg/post-fetch/domains/processors"
)

type PostApiProcessor struct {
	ctx             context.Context
	logger          adapters.Logger
	db              adapters.DB
	mu              *sync.RWMutex
	postApiClient   clients.PostApiClient
	taskRepoFactory repositories.TaskRepoFactory
	postRepoFactory repositories.PostRepoFactory
}

func NewPostApiProcessor(
	ctx context.Context,
	logger adapters.Logger,
	pool adapters.Pool,
	postApiClient clients.PostApiClient,
	taskRepoFactory repositories.TaskRepoFactory,
	postRepoFactory repositories.PostRepoFactory,
) (pdc.PostApiProcessor, error) {
	db, err := pool.Get(adapters.PostKind)
	if err != nil {
		return nil, err
	}

	return &PostApiProcessor{
		ctx:             ctx,
		logger:          logger,
		db:              db,
		mu:              &sync.RWMutex{},
		postApiClient:   postApiClient,
		taskRepoFactory: taskRepoFactory,
		postRepoFactory: postRepoFactory,
	}, nil
}

func (p *PostApiProcessor) Process(task *entities.Task) {
	for pageId := task.PageOffset + 1; pageId <= task.PageOffset+task.PageLimit; pageId++ {
		pageId := pageId
		go func() {
			_ = p.fetchAndSavePosts(task.Id, pageId)
		}()
	}
}

func (p *PostApiProcessor) fetchAndSavePosts(taskId uint32, pageId uint32) error {
	logger := p.logger.WithCtx(p.ctx, "method", "fetch_and_save_posts", "page", pageId)

	p.mu.Lock()
	defer p.mu.Unlock()

	tx, err := p.db.Begin(p.ctx)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer tx.Rollback()

	taskRepo := p.taskRepoFactory.Create(p.ctx, tx)
	task, err := taskRepo.FindById(taskId)
	if err != nil {
		logger.Error(err)
		return err
	}

	resp, err := p.postApiClient.Fetch(p.ctx, pageId)
	if err != nil {
		task.Fail++
		logger.Error(err)
		return err
	}

	var posts entities.Posts
	for _, data := range resp.Data {
		posts = append(posts, &entities.Post{
			PostId: data.Id,
			UserId: data.UserId,
			Title:  data.Title,
			Body:   data.Body,
		})
	}

	err = p.postRepoFactory.Create(p.ctx, tx).SaveAll(posts)
	if err != nil {
		task.Fail++
		logger.Error(err)
		return err
	}

	task.Success++
	err = taskRepo.Save(task)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Infof("posts are fetched successfully for page %d", pageId)
	return nil
}
