package processors

import (
	"tiu-85/gorest-iman/pkg/common/domain/entities"
)

type PostApiProcessor interface {
	Process(task *entities.Task)
}
