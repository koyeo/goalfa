package service

import (
	"context"
	"github.com/shopspring/decimal"
)

type FooService interface {
	Ping(ctx context.Context) (out string, err error)
	Ping2(ctx context.Context) (out decimal.Decimal, err error)
	AddPost(ctx context.Context, in Post) (out Post, err error)
	QueryPost(ctx context.Context, in QueryPostIn) (out []Post, err error)
	TestGetArray(ctx context.Context) (out [][]string, err error)
	TestPostArray(ctx context.Context) (out [][]Post, err error)
}

type CID struct {
	value string
}

func (p CID) String() string {
	return p.value
}

type Post struct {
	Title   string   `json:"title,omitempty" label:"标题" validator:"required" description:""`
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
