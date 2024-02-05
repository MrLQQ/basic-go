package channel

import (
	"testing"
	"time"
)

/*
* Channel是Go里面用于并发的一个关键内置类型
*
* Channel的基本操作包括：
* * 声明：var ch chan T 这种形态，T是channel中你准备放的数据的类型。
* * 创建：make(chan T) 和 make(chan T,size) 两种，后者带了容量参数。
* * 发送数据到channel里面：ch<-data,用的是箭头。
* * 从channel里面读取数据：var:= <-ch。
* * 使用close来关闭channel：close(ch)。
 */
func TestChannel(t *testing.T) {
	// 声明
	//var ch chan struct{}
	// 声明并创建
	//ch1 := make(chan int)
	// 这种是带buffer的
	ch2 := make(chan int, 3)
	// 把123发送到ch2里面
	ch2 <- 123
	data := <-ch2
	t.Log(data)
	// 这个是关闭channel
	close(ch2)
}

func TestChannelClose(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 0
	val, ok := <-ch
	if ok {
		t.Log("读取到了数据", val)
	}
	close(ch)
	// 这个操作会引起panic
	//ch <- 123
	val, ok = <-ch
	t.Log("读取到数据了吗？", ok, val)
}

func TestChannelLoop(t *testing.T) {
	ch := make(chan int, 1)
	go func() {
		for i := 0; i < 3; i++ {
			ch <- i
			time.Sleep(time.Second)
		}
		close(ch)
	}()
	for val := range ch {
		t.Log(val)
	}
}

func TestChannelBlocking(t *testing.T) {
	ch := make(chan int) // 注意该处，没有buffer，发送数据会阻塞
	b1 := BigStruct{}
	go func() {
		var b BigStruct
		// 在这里尝试发送数据，这个就是goroutine泄露
		ch <- 123
		t.Log(b, b1)
	}()
}

type BigStruct struct {
}

func TestChannelSelect(t *testing.T) {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 2)
	go func() {
		time.Sleep(time.Second * 2)
		ch1 <- 123
	}()
	go func() {
		time.Sleep(time.Second)
		ch2 <- 123
	}()
	select {
	case val := <-ch1:
		t.Log("进来了ch1", val)

	case val := <-ch2:
		t.Log("进来了ch2", val)
	}

}
