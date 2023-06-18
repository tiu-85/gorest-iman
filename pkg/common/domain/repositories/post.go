package repositories

import (
	"context"

	"tiu-85/gorest-iman/pkg/common/domain/entities"

	"tiu-85/gorest-iman/pkg/common/gen/pbv1"
	"tiu-85/gorest-iman/pkg/common/infra/adapters"
)

type PostRepoFactory interface {
	Create(context.Context, adapters.DB) PostRepo
}

type PostRepo interface {
	FindAllByFilter(*pbv1.GetPostListFilter, *pbv1.Paginator) (entities.Posts, uint32, error)
	FindById(uint32) (*entities.Post, error)
	Save(*entities.Post) error
	SaveAll(list entities.Posts) error
	Delete(post *entities.Post) error
}
