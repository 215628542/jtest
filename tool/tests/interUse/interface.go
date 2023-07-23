package interUse

import "context"

/*
	CourseClient结构体是课程顶层的抽象类，但是不直接定义成抽象类型
    是因为继承了这个课程顶层类型的话，就要实现这个抽象类的所有方法
    但是使用类中其实是不需要使用到那么多的方法

    - 注意：这个interface文件可以废弃掉
*/

//func NewCourseClient() *CourseClient {
//	return &CourseClient{}
//}

type CourseClient struct {
}

func (c *CourseClient) LearnGo(ctx context.Context, req interface{}) (resp interface{}, err error) {
	return
}

func (c *CourseClient) LearnJava(ctx context.Context, req interface{}) (resp interface{}, err error) {
	return
}

func (c *CourseClient) LearnPhp(ctx context.Context, req interface{}) (resp interface{}, err error) {
	return
}
