package main

import "fmt"

func main() {
	s := []int{1, 2, 3, 4, 5}
	fmt.Printf("原数组 len=%d,cap=%d,s=%v \n", len(s), cap(s), s)
	s = sliceDeleteIdx(s, 2)
	fmt.Printf("删除且缩容后的数组 len=%d,cap=%d,s=%v \n", len(s), cap(s), s)
}
