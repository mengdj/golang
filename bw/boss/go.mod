module boss

go 1.15

require (
	admin v0.0.0
	broadcast v0.0.0
	ext v0.0.0
	github.com/ThreeDotsLabs/watermill v1.1.1
	github.com/alecthomas/log4go v0.0.0-20180109082532-d146e6b86faa
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/panjf2000/ants/v2 v2.4.3
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b // indirect
	golang.org/x/sys v0.0.0-20201119102817-f84b799fce68 // indirect
	golang.org/x/text v0.3.4 // indirect
	golang.org/x/tools v0.0.0-20200117065230-39095c1d176c // indirect
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	server v0.0.0
	tool v0.0.0

)

replace (
	admin v0.0.0 => ./admin
	broadcast v0.0.0 => ./broadcast
	codec v0.0.0 => ./codec
	ext v0.0.0 => ./ext
	proc v0.0.0 => ./proc
	server v0.0.0 => ./server
	tool v0.0.0 => ./tool
)
