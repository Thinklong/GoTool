package main

import (
	"fmt"
	"time"
	//"runtime"
)

func say(s string) {
	for i := 0; i < 5; i++ {

		fmt.Println(s)
		//runtime.Gosched()
	}

}
func involution(c chan int64) {
	//fmt.Println(<-c)
	c <- time.Now().UnixNano()
	defer close(c)
	for i := 0; i < 100; i++ {
		//i := <-c
		//inv := i * i
		//c <- inv
		//fmt.Print("ceshi")
		c <- time.Now().UnixNano()
	}
}
func print_inv(c chan int64) {
	var v int64
	ok := true
	for ok {
		if v, ok = <-c; ok {
			fmt.Println(v)
		}
	}
}

func main() {
	//go say("word")
	//say("hello")
	//runtime.Gosched()
	//fmt.Println("hello word!!")
	c := make(chan int64)
	//c <- 2

	go involution(c)
	print_inv(c)
	//c <- 3
	//go involution(c)
	//c <- 43
	//print_inv(c)

	//fmt.Println(<-c)
}
