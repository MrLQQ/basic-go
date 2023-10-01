package main

import "fmt"

func main() {
	array := [5]int{1, 2, 3, 4, 5}
	s := array[:]
	fmt.Printf("原数组 \n")
	fmt.Printf("address of slice %p len=%d,cap=%d,s=%v \n", &array, len(array), cap(array), array)
	fmt.Printf("address of slice %p len=%d,cap=%d,s=%v \n", &s, len(s), cap(s), s)
	s = sliceDeleteIdx(s, 2)
	fmt.Printf("调用完成后的数组 \n")
	// 这里发现会修改原数组内容
	// 如果不想修改原数组内容，只想修改切片
	// 那在调用方法的内部，还是得再开辟一次新空间进行深拷贝后，再对新切片进行处理
	fmt.Printf("address of slice %p len=%d,cap=%d,s=%v \n", &array, len(array), cap(array), array)
	fmt.Printf("address of slice %p len=%d,cap=%d,s=%v \n", &s, len(s), cap(s), s)

}
