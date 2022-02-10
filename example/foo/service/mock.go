package service

import "context"

type FooMock struct {
}

func (f FooMock) Ping(ctx context.Context) (out string, err error) {
	out = "ok"
	return
}

func (f FooMock) AddPost(ctx context.Context, in Post) (out Post, err error) {
	out = in
	return
}

func (f FooMock) QueryPost(ctx context.Context, in QueryPostIn) (out []Post, err error) {
	out = append(out, Post{
		Title:   "一篇文章",
		Content: "文章内容",
		Tags:    []string{"A", "B", "C"},
	})
	return
}
