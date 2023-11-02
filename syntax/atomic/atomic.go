package main

import "sync/atomic"

func main() {
	var val int32 = 32
	// 原子读，你不会读到修改到一半的数据
	val = atomic.LoadInt32(&val)
	println(val)
	// 原子写，即便不同的Goroutine在不同的CPU核上，也可以立即看到。这个可以确保，缓存里面的val被改成14了
	atomic.StoreInt32(&val, 14)
	// 原子自增，返回自增后的结果
	newVal := atomic.AddInt32(&val, 1)
	println(newVal)
	// CAS错误
	// 如果val的值是14，就修改为15
	swapped := atomic.CompareAndSwapInt32(&val, 13, 15)
	println(swapped)
}
