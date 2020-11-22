module boss

go 1.15

require (
	broadcast v0.0.0
	github.com/alecthomas/log4go v0.0.0-20180109082532-d146e6b86faa
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/panjf2000/ants/v2 v2.4.3
	server v0.0.0
	tool v0.0.0

)

replace (
	broadcast v0.0.0 => ./broadcast
	codec v0.0.0 => ./codec
	proc v0.0.0 => ./proc
	server v0.0.0 => ./server
	tool v0.0.0 => ./tool
)
