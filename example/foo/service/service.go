package service

import "context"

type FooService interface {
	Ping(ctx context.Context) (out string, err error)
	AddPost(ctx context.Context, in Post) (out Post, err error)
	QueryPost(ctx context.Context, in QueryPostIn) (out []Post, err error)
}

type Post struct {
	Title   string   `json:"title,omitempty" label:"标题" validator:"required"`
	Content string   `json:"content,omitempty" label:"内容"`
	Tags    []string `json:"tags,omitempty" label:"标签"`
}

type QueryPostIn struct {
	Page  int    `json:"page"`
	Limit string `json:"limit"`
	Order *Order `json:"order"`
}

type Order struct {
	Field     string
	Direction string
}
