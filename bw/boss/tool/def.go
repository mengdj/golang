package tool

import "time"

const (
	//jwt
	TOKENEXPIREDURATION = time.Hour * 2
	TOKENSECRECT        = "生活不仅有远方的诗，还有眼前的苟且"
	//web
	DEFAULT_ADMIN string = "admin"
	RANDOM_CHARS         = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefhijklmnopqrstuvwxyz1234567890"
	//msg and websocket (WORKER_ADMIN_CTX为组，分支消息用元数据来区分)
	WORKER_ADMIN_CTX          = "WORKER_ADMIN_CTX"
	WORKER_ADMIN_CTX_RESULT   = "WORKER_ADMIN_CTX_RESULT"
	WORKER_ADMIN_CTX_SUB_TYPE = "type"
	WORKER_UPDATE             = "WORKER_UPDATE"
	WORKER_QPS                = "WORKER_QPS"
	WORKER_ONLINE             = "WORKER_ONLINE"
	WORKER_OFFLINE            = "WORKER_OFFLINE"

	QUERY_WORKERS           = "QUERY_WORKERS"
	QUERY_PING              = "QUERY_PING"
	REFRESH_SELECTED_WORKER = "REFRESH_SELECTED_WORKER"
	//process后台部分慢查询
	PROCESS_FOR_SLOW      = "PROCESS_FOR_SLOW"
	PROCESS_FOR_SLOW_TYPE = "type"
)
