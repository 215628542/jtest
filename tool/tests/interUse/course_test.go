package interUse

import (
	"context"
	"fmt"
	"testing"
)

func Test_courseServiceImpl_LearnGo(t *testing.T) {
	mockService := NewCourseService(&mockCourseClient{})
	req := make(map[string]interface{}, 0)
	ctx := context.Background()
	mockClient, ok := mockService.(*CourseServiceImpl)
	fmt.Printf("%#v", mockService)
	if ok {
		fmt.Println("mock test")
		mockClient.coursClient.LearnGo(ctx, req)
	} else {
		fmt.Println("not mock test")
	}
}
