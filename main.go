package main

import (
	"fmt"
	"os"
	"runtime"

	"hello-go/greeting"
)

// version перезаписывается через -ldflags "-X main.version=..." во время сборки.
// Локально остаётся "dev".
var version = "dev"

func main() {
	fmt.Printf("hello-go version %s\n", version)
	fmt.Println("Hello from Go! 🐹")
	fmt.Printf("OS: %s\n", runtime.GOOS)
	fmt.Printf("Arch: %s\n", runtime.GOARCH)
	fmt.Println(greeting.Greet("GitHub"))
	fmt.Printf("Sum 1..10 = %d\n", greeting.SumRange(1, 10))

	if len(os.Args) > 1 {
		fmt.Println("Аргументы:")
		for i, arg := range os.Args[1:] {
			fmt.Printf("  %d: %s\n", i+1, arg)
		}
	}
}
