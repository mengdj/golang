module boss

go 1.15

require (
	broadcast v0.0.0
	github.com/alecthomas/log4go v0.0.0-20180109082532-d146e6b86faa
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/panjf2000/ants/v2 v2.4.3
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b // indirect
	golang.org/x/sys v0.0.0-20201119102817-f84b799fce68 // indirect
	golang.org/x/text v0.3.4 // indirect
	golang.org/x/tools v0.0.0-20200117065230-39095c1d176c // indirect
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	ext v0.0.0
	server v0.0.0
	admin v0.0.0

)

replace (
	ext v0.0.0 => ./ext
	broadcast v0.0.0 => ./broadcast
	codec v0.0.0 => ./codec
	proc v0.0.0 => ./proc
	server v0.0.0 => ./server
	tool v0.0.0 => ./tool
	admin v0.0.0 => ./admin
)
