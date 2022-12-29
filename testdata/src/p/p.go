package p

import "fmt"

func safeGo(call func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
			}
		}()
		call()
	}()
}

func safeGoroutine() {
	safeGo(func() {
		fmt.Println("go go go")
	})
}
