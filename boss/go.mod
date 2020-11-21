module boss

go 1.15

require (
	broadcast v0.0.0
	tool v0.0.0
	github.com/ThreeDotsLabs/watermill v1.1.1 // indirect
	github.com/alecthomas/log4go v0.0.0-20180109082532-d146e6b86faa
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/lrita/gosync v0.0.0-20180507073543-17c653898443 // indirect
	github.com/panjf2000/ants/v2 v2.4.3
	github.com/panjf2000/gnet v1.3.1 // indirect
	github.com/robfig/cron v1.2.0 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	proc v0.0.0
	codec v0.0.0
	server v0.0.0

)

replace (
	tool v0.0.0 => ./tool
	codec v0.0.0 => ./codec
	broadcast v0.0.0 => ./broadcast
	proc v0.0.0 => ./proc
	server v0.0.0 => ./server
)
