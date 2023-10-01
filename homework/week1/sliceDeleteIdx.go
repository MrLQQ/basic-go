package main

import (
	"fmt"
)

/*
实现删除切片特定下标元素的方法。

要求一：能够实现删除操作就可以。
要求二：考虑使用比较高性能的实现。
要求三：改造为泛型方法
要求四：支持缩容，并旦设计缩容机制。
*/
func sliceDeleteIdx[T any](s []T, idx int) []T {
	println("方法内部数组")
	fmt.Printf("address of slice %p len=%d,cap=%d,s=%v \n", &s, len(s), cap(s), s)

	// 校验
	if idx < 0 || idx > len(s) {
		panic("非法的idx")
	}
	// 方式1：通过位移删除目标元素
	//for i := idx; i < len(s)-1; i++ {
	//	s[i] = s[i+1]
	//}
	// 方式2：通过截取新切片删除目标元素
	s = append(s[:idx], s[idx+1:]...)
	fmt.Printf("address of slice %p len=%d,cap=%d,s=%v \n", &s, len(s), cap(s), s)

	// 开辟新空间
	// 如果要缩容的话，目前只能想到开辟一个新空间来深拷贝切片
	newSlice := make([]T, len(s))
	// 深拷贝
	copy(newSlice, s[:len(s)])
	fmt.Printf("address of slice %p len=%d,cap=%d,s=%v \n", &newSlice, len(newSlice), cap(newSlice), newSlice)
	return newSlice
}
