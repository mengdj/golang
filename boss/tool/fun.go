package tool

func Try(f func(),e func(i interface{})){
	defer func() {
		if r:=recover();nil!=e{
			e(r)
		}
	}()
	f()
}
