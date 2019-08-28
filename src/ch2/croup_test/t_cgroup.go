package main

import (
	"fmt"
	"os"
	"syscall"
)

func main()  {
	fmt.Printf("current pro id is %d", syscall.Getpid())
	if os.Args[0] == "/proc/self/exe" {
		// 容器进程
		fmt.Printf("current pro id is %d", syscall.Getpid())
	}

}
