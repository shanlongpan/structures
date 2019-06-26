package main

import (
	"fmt"
	"time"
)
//斐波那契数列非递归
func Fibonacci(n int)int{
	if n>1{
		f0:=0;f1:=1;f2:=0
		for i:=1;i<n;i++{
			f2=f0+f1
			f0=f1;f1=f2
		}
		return f2
	}
	return 1
}

//斐波那契数列递归
func fib(x int) int {
	if x < 2 {
		return x
	}
	return fib(x-1) + fib(x-2)
}

func main()  {
	now:=time.Now()
	fmt.Println(Fibonacci(40))
	fmt.Println(time.Since(now))
	now=time.Now()
	fmt.Println(fib(40))
	fmt.Println(time.Since(now))
}
