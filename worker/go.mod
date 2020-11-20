module worker

go 1.15

require (
	ext v0.0.0
	broadcast v0.0.0
	client v0.0.0
	github.com/ThreeDotsLabs/watermill v1.1.1
	github.com/alecthomas/log4go v0.0.0-20180109082532-d146e6b86faa
	github.com/golang/protobuf v1.4.3
	github.com/kbinani/screenshot v0.0.0-20191211154542-3a185f1ce18f
	github.com/lxn/win v0.0.0-20201111105847-2a20daff6a55 // indirect
	github.com/shirou/gopsutil v3.20.10+incompatible // indirect
	github.com/spf13/viper v1.7.1 // indirect
	proc v0.0.0
	screenshot v0.0.0

)

replace (
	ext v0.0.0 => ./ext
	broadcast v0.0.0 => ./broadcast
	client v0.0.0 => ./client
	codec v0.0.0 => ./codec
	proc v0.0.0 => ./proc
	screenshot v0.0.0 => ./screenshot
)
