module broadcast

go 1.15

require (
	github.com/ThreeDotsLabs/watermill v1.1.1
	github.com/alecthomas/log4go v0.0.0-20180109082532-d146e6b86faa
	github.com/asaskevich/EventBus v0.0.0-20200907212545-49d423059eef
	github.com/golang/protobuf v1.4.3
	github.com/panjf2000/ants/v2 v2.4.3
	github.com/panjf2000/gnet v1.3.1
	github.com/robfig/cron/v3 v3.0.1
	proc v0.0.0
	ext v0.0.0
)

replace proc v0.0.0 => ../proc
replace ext v0.0.0 => ../ext
