package main

/*
实现删除切片特定下标元素的方法。

要求一：能够实现删除操作就可以。
要求二：考虑使用比较高性能的实现。
要求三：改造为泛型方法
要求四：支持缩容，并旦设计缩容机制。
*/
func sliceDeleteIdx[T any](s []T, idx int) []T {
	// 校验
	if idx < 0 || idx > len(s) {
		panic("非法的idx")
	}
	// 位移
	for i := idx; i < len(s)-1; i++ {
		s[i] = s[i+1]
	}
	// 开辟新空间
	newSlice := make([]T, len(s)-1)
	// 深拷贝
	copy(newSlice, s[:len(s)-1])
	return newSlice
}
