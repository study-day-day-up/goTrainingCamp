package main

import "fmt"

func main() {
	n := 10
	fmt.Println(Fibonacci(n))
}

func Fibonacci(n int) int {
	a, b := 0, 1
	fmt.Printf("Fibonacci(%d): %d\n", a, a)

	for i := 1; i <= n; i++ {
		a, b = b, a+b
		//依次输出 Fibonacci 数列
		fmt.Printf("Fibonacci(%d): %d\n", i, a)
	}
	return a
}
