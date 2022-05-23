package main

import "fmt"

func main() {
	fmt.Println(Max(1, 4))
	fmt.Println(Max(5, 4))
}

func Max(x, y int) int {
	return Ternary(x > y, x, y)
}

func Ternary[T any](cond bool, x T, y T) T {
	if cond {
		return x
	}
	return y
}
