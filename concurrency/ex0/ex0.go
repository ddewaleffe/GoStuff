package main

import (
	"fmt"
	"math/cmplx"
	"time"
)

func writer(c chan complex128, w chan bool) {
	for i := 0; i < 10; i++ {
		v:= complex(float64(i), float64(i*7))
		fmt.Printf("Sending %v\n",v);
		c <- v
	}
	//w <- true
}
func reader(c chan complex128, w chan bool) {
	for {
		select {
		case  v := <-c :
			fmt.Printf("Got %v, sqrt=>%v\n", v, cmplx.Sqrt(v))
		case  <- time.After(8000 * time.Millisecond):
			fmt.Printf("Timout. Quitting.\n")
			w<-true
			}
	}
	//w<- true
}
func main() {

	w := make(chan bool)
	c := make(chan complex128)
	go reader(c, w)
	writer(c, w)

	//wait for ending signal...
	<-w

}
