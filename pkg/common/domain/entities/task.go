package entities

import (
	"tiu-85/gorest-iman/pkg/common/gen/pbv1"
)

type Task struct {
	tableName struct{} `pg:"post.tasks"`

	Id         uint32
	Total      uint32 `pg:",use_zero"`
	Success    uint32 `pg:",use_zero"`
	Fail       uint32 `pg:",use_zero"`
	PageLimit  uint32 `pg:",use_zero"`
	PageOffset uint32 `pg:",use_zero"`
}

func (p *Task) ToProto() *pbv1.Task {
	return &pbv1.Task{
		Id:         p.Id,
		Total:      p.Total,
		Success:    p.Success,
		Fail:       p.Fail,
		PageLimit:  p.PageLimit,
		PageOffset: p.PageOffset,
	}
}

type Tasks []*Task

func (p Tasks) ToProto() []*pbv1.Task {
	list := make([]*pbv1.Task, len(p))
	for i, item := range p {
		list[i] = item.ToProto()
	}

	return list
}
