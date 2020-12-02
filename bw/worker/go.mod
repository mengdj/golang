module worker

go 1.15

require (
	broadcast v0.0.0
	client v0.0.0
	ext v0.0.0
	github.com/ThreeDotsLabs/watermill v1.1.1
	github.com/alecthomas/log4go v0.0.0-20180109082532-d146e6b86faa
	github.com/golang/protobuf v1.4.3
	proc v0.0.0
	screenshot v0.0.0
)

replace (
	broadcast v0.0.0 => ./broadcast
	client v0.0.0 => ./client
	codec v0.0.0 => ./codec
	ext v0.0.0 => ./ext
	proc v0.0.0 => ./proc
	screenshot v0.0.0 => ./screenshot
	tool v0.0.0 => ./tool
)
