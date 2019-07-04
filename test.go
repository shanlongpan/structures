package main

import (
	"bytes"
	"context"
	"fmt"
	"time"
)
var shardIndexes = make([]uint64, 32)

func main() {

	s := []byte("同学们，上午好")
	m := func(r rune) rune {
		if r == '上' {
			r = '下'
		}
		return r
	}

	fmt.Println(string(s))
	fmt.Println(string(bytes.Map(m, s)))

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	//defer cancel()

	go handle(ctx, 1500*time.Millisecond)
	cancel()
	select {
	case <-ctx.Done():
		fmt.Println("main", ctx.Err())
	}
}

func handle(ctx context.Context, duration time.Duration) {
	select {
	case <-ctx.Done():
		fmt.Println("handle", ctx.Err())

	case <-time.After(duration):
		fmt.Println("process request with", duration)
	}
}
