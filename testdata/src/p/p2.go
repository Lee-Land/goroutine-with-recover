package p

import "fmt"

func rawGoroutine() {
	go func() {
		fmt.Println("go go go")
	}()
}
