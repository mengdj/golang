module admin

go 1.15

require (
	ext v0.0.0
	github.com/ThreeDotsLabs/watermill v1.1.1
	github.com/alecthomas/log4go v0.0.0-20180109082532-d146e6b86faa
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.7.0
	github.com/go-session/gin-session v3.1.0+incompatible
	github.com/go-session/session v3.1.2+incompatible // indirect
	github.com/gogf/gf v1.14.5
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/json-iterator/go v1.1.10
	github.com/panjf2000/ants/v2 v2.4.3
	github.com/smartystreets/goconvey v1.6.4 // indirect
	gopkg.in/olahol/melody.v1 v1.0.0-20170518105555-d52139073376
	tool v0.0.0
)

replace tool v0.0.0 => ../tool

replace ext v0.0.0 => ../ext
