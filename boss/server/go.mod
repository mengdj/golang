module server

go 1.15

require (
	github.com/alecthomas/log4go v0.0.0-20180109082532-d146e6b86faa
	github.com/golang/protobuf v1.4.3
	github.com/panjf2000/gnet v1.3.1
	github.com/robfig/cron v1.2.0
	codec v0.0.0
	proc v0.0.0
)

replace (
	proc v0.0.0 => ../proc
	codec v0.0.0 => ../codec
)
