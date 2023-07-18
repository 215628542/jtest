package taskFunc

import (
	"encoding/json"
	"errors"
	"fmt"
	"test/tool/asnycqueue/core"
)

type TaskFunc struct {
}

type TestQueueMsg struct {
	TestMsg string `json:"test_msg"`
}

func Test(data string) error {
	//panic(123456)

	reqData := core.AsyncReqData{}
	json.Unmarshal([]byte(data), &reqData)

	msg := TestQueueMsg{}
	json.Unmarshal([]byte(reqData.Value), &msg)

	fmt.Printf("异步任务获取数据--%s\n", msg.TestMsg)
	return errors.New("test error=========")
}

// @stackInfo string panic堆返回的异常信息
func TestPanicFunc(reqData string, stackInfo []byte) {
	fmt.Println("TestPanicFunc获取请求信息-" + reqData)
	fmt.Println("TestPanicFunc获取堆栈信息" + string(stackInfo))
}
