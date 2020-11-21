module server

go 1.15

require (
	codec v0.0.0
	tool v0.0.0
	github.com/alecthomas/log4go v0.0.0-20180109082532-d146e6b86faa
	github.com/gogf/gf v1.14.4
	github.com/golang/protobuf v1.4.3
	github.com/panjf2000/gnet v1.3.1
	github.com/robfig/cron v1.2.0
	proc v0.0.0
)

replace (
	tool v0.0.0 => ../tool
	codec v0.0.0 => ../codec
	proc v0.0.0 => ../proc
)
