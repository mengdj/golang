package main_test

import (
	"fmt"
	"testing"
	"tool"
)
/** 测试获取本机IP地址 */
func TestLocalNetAddress(t *testing.T) {
	fmt.Println(tool.LocalAddress())
}
