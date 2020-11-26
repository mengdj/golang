module broadcast

go 1.15

require (
	github.com/alecthomas/log4go v0.0.0-20180109082532-d146e6b86faa
	github.com/golang/protobuf v1.4.3
	github.com/robfig/cron/v3 v3.0.1
	tool v0.0.0
	proc v0.0.0
)

replace proc v0.0.0 => ../proc
replace tool v0.0.0 => ../tool