package interUse

import (
	"context"
	"fmt"
)

type mockCourseClient struct {
}

func (c *mockCourseClient) LearnGo(ctx context.Context, req interface{}) (resp interface{}, err error) {

	fmt.Println("mock学习golang方法")
	return
}
