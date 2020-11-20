module client

go 1.15

require (
	ext v0.0.0
	codec v0.0.0
	github.com/BurntSushi/xgb v0.0.0-20201008132610-5f9e7b3c49cd // indirect
	github.com/ThreeDotsLabs/watermill v1.1.1
	github.com/alecthomas/log4go v0.0.0-20180109082532-d146e6b86faa
	github.com/gen2brain/shm v0.0.0-20200228170931-49f9650110c5 // indirect
	github.com/kbinani/screenshot v0.0.0-20191211154542-3a185f1ce18f
	github.com/lrita/gosync v0.0.0-20180507073543-17c653898443
	github.com/robfig/cron v1.2.0
	proc v0.0.0
	screenshot v0.0.0
)

replace (
	ext v0.0.0 => ../ext
	codec v0.0.0 => ../codec
	proc v0.0.0 => ../proc
	screenshot v0.0.0 => ../screenshot
)
