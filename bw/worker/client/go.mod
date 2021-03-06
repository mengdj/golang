module client

go 1.15

require (
	codec v0.0.0
	ext v0.0.0
	github.com/ThreeDotsLabs/watermill v1.1.1
	github.com/alecthomas/log4go v0.0.0-20180109082532-d146e6b86faa
	github.com/gogf/gf v1.14.5
	github.com/golang/protobuf v1.4.3
	github.com/golang/snappy v0.0.2
	github.com/robfig/cron v1.2.0
	proc v0.0.0
	screenshot v0.0.0
)

replace (
	codec v0.0.0 => ../codec
	ext v0.0.0 => ../ext
	proc v0.0.0 => ../proc
	screenshot v0.0.0 => ../screenshot
	tool v0.0.0 => ../tool
)
