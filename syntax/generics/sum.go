package main

import (
	"io"
)

func Sum[T Number](vals []T) T {
	var res T
	for _, v := range vals {
		res = res + v
	}
	return res
}

func Max[T Number](vars []T) T {
	t := vars[0]
	for i := 1; i < len(vars); i++ {
		if t < vars[i] {
			t = vars[i]
		}
	}
	return t
}

func Find[T any](vals []T, filter func(T) bool) T {
	for _, v := range vals {
		if filter(v) {
			return v
		}
	}
	var t T
	return t
}

func Inster[T any](idx int, val T, vals []T) []T {
	if idx < 0 || idx >= len(vals) {
		panic("idx不合法")
	}
	// 先扩容
	vals = append(vals, val)
	for i := len(vals) - 2; i >= idx; i-- {
		vals[i], vals[i+1] = vals[i+1], vals[i]
	}
	return vals
}

type Integer int

type Number interface {
	~int | uint | int32
}

func UseSum() {
	res := Sum[int]([]int{123, 123})
	println(res)
	resV1 := Sum[Integer]([]Integer{123, 123})
	println(resV1)
}

func Closable[T io.Closer]() {
	var t T
	t.Close()
}
