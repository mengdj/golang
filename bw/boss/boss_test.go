package main_test

import (
	"fmt"
	"github.com/reactivex/rxgo/v2"
	"testing"
	"tool"
)

/** 测试获取本机IP地址 */
func TestLocalNetAddress(t *testing.T) {
	fmt.Println(tool.LocalAddress())
}

func TestWorkFlow(t *testing.T) {
	observable := rxgo.Just(1, 2, 3)()
	observable.Filter(func(i interface{}) bool {
		x:=i.(int)
		if x%2!=0{
			return true
		}
		return false
	}).DoOnNext(func(i interface{}) {
		fmt.Println(i)
	})
}
