package entities

import (
	"tiu-85/gorest-iman/pkg/common/gen/pbv1"
)

type Post struct {
	tableName struct{} `pg:"post.posts"`

	Id     uint32
	PostId uint32
	UserId uint32
	Title  string
	Body   string
}

func (p *Post) ToProto() *pbv1.Post {
	return &pbv1.Post{
		Id:     p.Id,
		PostId: p.PostId,
		UserId: p.UserId,
		Title:  p.Title,
		Body:   p.Body,
	}
}

type Posts []*Post

func (p Posts) ToProto() []*pbv1.Post {
	list := make([]*pbv1.Post, len(p))
	for i, item := range p {
		list[i] = item.ToProto()
	}

	return list
}
