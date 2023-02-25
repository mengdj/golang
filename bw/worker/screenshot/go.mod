module screenshot

go 1.15

require (
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/shirou/gopsutil v3.20.10+incompatible
	github.com/spf13/viper v1.7.1
	golang.org/x/sys v0.1.0 // indirect
	tool v0.0.0
)

replace tool v0.0.0 => ../tool
