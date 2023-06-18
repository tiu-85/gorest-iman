package services

import (
	"context"

	"tiu-85/gorest-iman/pkg/common/domain/repositories"
	"tiu-85/gorest-iman/pkg/common/gen/pbv1"
	"tiu-85/gorest-iman/pkg/common/infra/adapters"
)

type PostCrudService struct {
	logger          adapters.Logger
	db              adapters.DB
	postRepoFactory repositories.PostRepoFactory
}

func NewPostCrudService(
	logger adapters.Logger,
	pool adapters.Pool,
	postRepoFactory repositories.PostRepoFactory,
) (*PostCrudService, error) {
	db, err := pool.Get(adapters.PostKind)
	if err != nil {
		return nil, err
	}

	return &PostCrudService{
		logger:          logger,
		db:              db,
		postRepoFactory: postRepoFactory,
	}, nil
}

func (s *PostCrudService) GetPostList(ctx context.Context, req *pbv1.GetPostListRequest) (*pbv1.GetPostListResponse, error) {
	logger := s.logger.WithCtx(ctx, "method", "get_post_list")
	postRepo := s.postRepoFactory.Create(ctx, s.db)

	posts, count, err := postRepo.FindAllByFilter(req.Filter, req.Paginator)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &pbv1.GetPostListResponse{
		Items: posts.ToProto(),
		Count: count,
	}, nil
}

func (s *PostCrudService) GetPost(ctx context.Context, req *pbv1.GetPostRequest) (*pbv1.GetPostResponse, error) {
	logger := s.logger.WithCtx(ctx, "method", "get_post")
	postRepo := s.postRepoFactory.Create(ctx, s.db)

	post, err := postRepo.FindById(req.Id)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &pbv1.GetPostResponse{
		Post: post.ToProto(),
	}, nil
}

func (s *PostCrudService) EditPost(ctx context.Context, req *pbv1.EditPostRequest) (*pbv1.EditPostResponse, error) {
	logger := s.logger.WithCtx(ctx, "method", "edit_post")
	postRepo := s.postRepoFactory.Create(ctx, s.db)

	post, err := postRepo.FindById(req.Id)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	post.Title = req.Params.Title
	post.Body = req.Params.Body

	err = postRepo.Save(post)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &pbv1.EditPostResponse{
		Post: post.ToProto(),
	}, nil
}

func (s *PostCrudService) DeletePost(ctx context.Context, req *pbv1.DeletePostRequest) (*pbv1.DeletePostResponse, error) {
	logger := s.logger.WithCtx(ctx, "method", "delete_post")
	postRepo := s.postRepoFactory.Create(ctx, s.db)

	post, err := postRepo.FindById(req.Id)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	err = postRepo.Delete(post)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &pbv1.DeletePostResponse{}, nil
}
