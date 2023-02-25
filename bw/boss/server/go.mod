module server

go 1.15

require (
	admin v0.0.0
	codec v0.0.0
	ext v0.0.0
	github.com/ThreeDotsLabs/watermill v1.1.1
	github.com/alecthomas/log4go v0.0.0-20180109082532-d146e6b86faa
	github.com/gogf/gf v1.14.5
	github.com/golang/protobuf v1.4.3
	github.com/golang/snappy v0.0.2
	github.com/google/go-cmp v0.5.3 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/json-iterator/go v1.1.10
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/panjf2000/gnet v1.3.1
	github.com/pkg/errors v0.9.1 // indirect
	github.com/robfig/cron v1.2.0
	github.com/stretchr/testify v1.6.1 // indirect
	golang.org/x/net v0.7.0 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	proc v0.0.0
	tool v0.0.0
)

replace (
	admin v0.0.0 => ../admin
	codec v0.0.0 => ../codec
	ext v0.0.0 => ../ext
	proc v0.0.0 => ../proc
	tool v0.0.0 => ../tool
)
