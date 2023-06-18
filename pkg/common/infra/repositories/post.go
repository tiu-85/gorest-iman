package repositories

import (
	"context"
	"errors"

	"github.com/go-pg/pg/v9"

	"tiu-85/gorest-iman/pkg/common/domain/entities"

	"tiu-85/gorest-iman/pkg/common/domain/repositories"
	"tiu-85/gorest-iman/pkg/common/gen/pbv1"
	"tiu-85/gorest-iman/pkg/common/infra/adapters"
	"tiu-85/gorest-iman/pkg/common/infra/values"
)

func NewPostRepoFactory(
	ctx context.Context,
	logger adapters.Logger,
) repositories.PostRepoFactory {
	return &postRepoFactory{
		ctx:    ctx,
		logger: logger,
	}
}

type postRepoFactory struct {
	ctx    context.Context
	logger adapters.Logger
}

func (f *postRepoFactory) Create(ctx context.Context, db adapters.DB) repositories.PostRepo {
	return &postRepo{
		DS:     adapters.NewDS(db.WithCtx(ctx)),
		logger: f.logger.WithCtx(ctx, "repo", "post"),
	}
}

type postRepo struct {
	*adapters.DS
	logger adapters.Logger
}

func (r *postRepo) FindById(id uint32) (*entities.Post, error) {
	logger := r.logger.With("method", "find_by_id")
	model := new(entities.Post)

	err := r.Model(model).Where("post.id = ?", id).Select()
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

func (r *postRepo) Save(model *entities.Post) error {
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

func (r *postRepo) SaveAll(list entities.Posts) error {
	insertList := entities.Posts{}
	updateList := entities.Posts{}

	for _, item := range list {
		if item.Id == 0 {
			insertList = append(insertList, item)
		} else {
			updateList = append(updateList, item)
		}
	}

	if len(insertList) > 0 {
		if _, err := r.Model(&insertList).Insert(&insertList); err != nil {
			r.logger.With("method", "save_all").Error(err)
			return values.ErrCantSave
		}
	}

	if len(updateList) > 0 {
		if _, err := r.Model(&updateList).Update(&updateList); err != nil {
			r.logger.With("method", "save_all").Error(err)
			return values.ErrCantSave
		}
	}

	return nil
}

func (r *postRepo) Delete(model *entities.Post) error {
	logger := r.logger.With("method", "delete")

	res, err := r.Model(model).Where("id = ?", model.Id).Delete()
	if err != nil {
		logger.Error(err)
		return values.ErrCantDelete
	}

	if res.RowsAffected() == 0 {
		logger.Error("no rows affected")
		return values.ErrCantDelete
	}

	return nil
}

func (r *postRepo) FindAllByFilter(filter *pbv1.GetPostListFilter, paginator *pbv1.Paginator) (entities.Posts, uint32, error) {
	logger := r.logger.With("method", "find_all_by_filter")

	var model entities.Posts
	query := r.Model(&model)

	if paginator != nil {
		query.Limit(int(paginator.Limit)).Offset(int(paginator.Offset))
	}

	if filter != nil && filter.PostId > 0 {
		query.Where("post.post_id = ?", filter.PostId)
	}

	if filter != nil && filter.UserId > 0 {
		query.Where("post.user_id = ?", filter.UserId)
	}

	if filter != nil && filter.Title != "" {
		query.Where("post.title like ?", filter.Title)
	}

	count, err := query.Count()
	if err != nil {
		logger.Error(err)
		return nil, 0, err
	}

	err = query.Select()
	if err != nil {
		logger.Error(err)
		return nil, 0, err
	}

	return model, uint32(count), nil
}
