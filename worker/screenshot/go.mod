module screenshot

go 1.15

require (
	tool v0.0.0
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/shirou/gopsutil v3.20.10+incompatible
	github.com/spf13/viper v1.7.1
	golang.org/x/sys v0.0.0-20201119102817-f84b799fce68 // indirect
)

replace (
	tool v0.0.0 => ../tool
)
