package clients

import (
	"context"
)

type PostApiClient interface {
	Fetch(context.Context, uint32) (*PostFetchResponse, error)
}

type PostApiClientRequest struct {
	Method  string
	Url     string
	Path    string
	Headers map[string]string
}

type PostFetchResponseData struct {
	Id     uint32 `json:"id"`
	UserId uint32 `json:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type PostFetchResponse struct {
	Data []PostFetchResponseData
}
