package tool

//尝试调用某个可能panic的函数，如果出错就捕捉，避免退出
func Try(f func(),e func(i interface{})){
	defer func() {
		if r:=recover();nil!=r{
			e(r)
		}
	}()
	f()
}