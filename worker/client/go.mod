module client

go 1.15

require (
	codec v0.0.0
	ext v0.0.0
	github.com/BurntSushi/xgb v0.0.0-20201008132610-5f9e7b3c49cd // indirect
	github.com/ThreeDotsLabs/watermill v1.1.1
	github.com/alecthomas/log4go v0.0.0-20180109082532-d146e6b86faa
	github.com/gen2brain/shm v0.0.0-20200228170931-49f9650110c5 // indirect
	github.com/gogf/gf v1.14.4
	github.com/golang/protobuf v1.4.3
	github.com/golang/snappy v0.0.2
	github.com/kbinani/screenshot v0.0.0-20191211154542-3a185f1ce18f
	github.com/lrita/gosync v0.0.0-20180507073543-17c653898443
	github.com/robfig/cron v1.2.0
	proc v0.0.0
	screenshot v0.0.0
)

replace (
	codec v0.0.0 => ../codec
	ext v0.0.0 => ../ext
	proc v0.0.0 => ../proc
	screenshot v0.0.0 => ../screenshot
)
