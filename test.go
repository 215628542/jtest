package main

import (
	"fmt"
	"runtime"
	"strings"
)

type T func([]int, int)

func T1(sl []int, i int) T {

	return func(sl []int, i int) {

		defer func() {
			if x := recover(); x != nil {
				_, file, line, ok := runtime.Caller(1) // skip the first frame (panic itself)
				if ok && strings.Contains(file, "runtime/") {
					// The panic came from the runtime, most likely due to incorrect
					// map/slice usage. The parent frame should have the real trigger.
					_, file, line, ok = runtime.Caller(2)
				}

				// Include the file and line number info in the error, if runtime.Caller returned ok.
				if ok {
					fmt.Printf("panic--- [%s:%d]: %v\n", file, line, x)
				} else {
					fmt.Printf("panic----: %v\n", x)
				}
				fmt.Println(x)
			}
		}()
		fmt.Println(sl[i])
	}
}

func T2(sl []int, i int) T {

	return func(sl []int, i int) {

		defer func() {
			if x := recover(); x != nil {
				_, file, line, ok := runtime.Caller(1) // skip the first frame (panic itself)
				if ok && strings.Contains(file, "runtime/") {
					// The panic came from the runtime, most likely due to incorrect
					// map/slice usage. The parent frame should have the real trigger.
					_, file, line, ok = runtime.Caller(2)
				}

				// Include the file and line number info in the error, if runtime.Caller returned ok.
				if ok {
					fmt.Printf("panic [%s:%d]: %v\n", file, line, x)
				} else {
					fmt.Printf("panic: %v\n", x)
				}
			}
		}()
		fmt.Println(sl[i])
	}
}

func main2() {

	a := []int{1}
	b := []int{2}

	funcMap := make([]T, 0)
	var t1, t2 T
	t1 = T1(a, 1)
	t2 = T2(b, 1)

	funcMap = append(funcMap, t1, t2)
	for _, f := range funcMap {
		f([]int{}, 1)
	}

	fmt.Println("end====")

}
