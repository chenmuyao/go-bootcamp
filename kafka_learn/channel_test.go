package kafkalearn

import (
	"testing"
	"time"
)

func TestChannel(t *testing.T) {
	// var ch chan struct{}
	// ch1 := make(chan int)
	ch2 := make(chan int, 3)
	ch2 <- 123
	data := <-ch2
	t.Log(data)
	close(ch2)
}

func TestChannelClose(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 0
	val, ok := <-ch
	t.Log("read?", ok, val)
	close(ch)
	// panic
	// ch <- 123
	val, ok = <-ch
	t.Log("read?", ok, val)

	// close again will panic
	// NOTE: the creator of the channel should close the channel
}

func TestForLoop(t *testing.T) {
	ch := make(chan int, 1)
	go func() {
		for i := range 3 {
			ch <- i
			time.Sleep(time.Millisecond)
		}
		close(ch)
	}()
	for val := range ch {
		t.Log(val)
	}
}

func TestChannelBlocking(t *testing.T) {
	ch := make(chan int)
	go func() {
		// leak
		ch <- 123
	}()

	ch1 := make(chan int)
	bigStruct1 := struct{}{}
	go func() {
		bigStruct := struct{}{}
		ch1 <- 123
		// leak
		t.Log(bigStruct, bigStruct1)
	}()
}

func TestSelect(t *testing.T) {
	for range 10 {
		ch1 := make(chan int, 1)
		ch2 := make(chan int, 2)
		go func() {
			time.Sleep(time.Millisecond)
			ch1 <- 123
		}()
		go func() {
			time.Sleep(time.Millisecond)
			ch2 <- 123
		}()
		select {
		case val := <-ch1:
			t.Log("ch1", val)
		case val := <-ch2:
			t.Log("ch2", val)
		}
	}
}
