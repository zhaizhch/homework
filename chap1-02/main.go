package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 10)
	done := make(chan bool)
	count := 1
	//consumer
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for _ = range ticker.C {
			select {
			case <-done:
				fmt.Println("consumer process interrupt...")
				return
			default:
				ch <- count
				fmt.Println("ch <-", count)
				count++
				fmt.Println("ch's size is ", len(ch))
			}
		}
	}()
	//producer
	go func() {
		time.Sleep(time.Second * 15) //判断10s后channel是否会堵塞
		ticker := time.NewTicker(1 * time.Second)
		for _ = range ticker.C {
			select {
			case <-done:
				fmt.Println("producer process interrupt...")
				return
			case i := <-ch:
				fmt.Println("i = ", i)
			}
		}
	}()
	//20s后done关闭
	time.Sleep(time.Second * 20)
	close(done)
	time.Sleep(time.Second)
	println("main function exit")
}
