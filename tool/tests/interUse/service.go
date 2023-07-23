package interUse

import (
	"context"
	"fmt"
)

/*
	使用方
	其实是将LearnGo方法抽象出来，变成多态，方便单测mock时，切换LearnGo方法调用的地方
*/

type CourseService interface {
}

type CourseProxy interface {
	LearnGo(ctx context.Context, req interface{}) (resp interface{}, err error)
}

type CourseServiceImpl struct {
	coursClient CourseProxy
}

func NewCourseService(courseClient CourseProxy) CourseService {
	return &CourseServiceImpl{
		coursClient: courseClient,
	}
}

func (c *CourseServiceImpl) LearnGo(ctx context.Context, req interface{}) (resp interface{}, err error) {

	fmt.Println("学习golang")
	return
}
