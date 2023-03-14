package types

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func process(ch chan string) {
	time.Sleep(10 * time.Second)
	ch <- "hi hello"
}

func TestContext(t *testing.T) {
	cancelCtx, cancel := context.WithCancel(context.Background())
	ctx := NewContext().WithContext(cancelCtx).WithCancel(cancel)
	ch := make(chan string)
	go process(ch)
	go func() {
		time.Sleep(5 * time.Second)
		ctx.Close()
	}()

	go func() {
		for range time.Tick(1000 * time.Millisecond) {
			fmt.Println("otherone")
		}
	}()

	for range time.Tick(1000 * time.Millisecond) {
		select {
		case v := <-ch:
			println(v)
			continue
		case <-ctx.Context().Done():
			fmt.Println("done context", ctx.Context().Err())
			return
		default:
			println("aaa")
		}

		println("aaaa")
	}
}
