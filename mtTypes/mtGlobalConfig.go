package mttypes

import (
	"os"
	"time"
)

// docker配置
var NODE_ID = os.Getenv("NODE_ID")

// redis配置
var (
	RedisConfig = RedisConf{
		Addr:    os.Getenv("REDIS_URL"),
		Pwd:     "reins5401",
		DB:      "0",
		Timeout: "10",
	}
)

// influxdb配置
var (
	InfluxConfig = InfluxConf{
		Url:    os.Getenv("INFLUXDB_URL"),
		Token:  os.Getenv("INFLUXDB_TOKEN"),
		Bucket: os.Getenv("INFLUXDB_BUCKET"),
		Org:    "miot-tracer",
	}
)

// log 配置
var LogAddr = "./logFile/miot_tracer_log.json"

//-------------------以下是测试配置-------------------

// // test配置
// var NODE_ID = "1"

// // redis配置
// var (
// 	RedisConfig = RedisConf{
// 		Addr:    "localhost:6379",
// 		Pwd:     "reins5401",
// 		DB:      "0",
// 		Timeout: "10",
// 	}
// )

// // influxdb配置
// var (
// 	InfluxConfig = InfluxConf{
// 		Url:    "http://localhost:8086",
// 		Token:  "J_xeoyLkPQFHBilXk4ELHjV85A7fFtIJvlo3GTjmKnF3QPZU63H7N0FH5_x7JBMPy3MRvVwoeoW0rnReDyLuPg==",
// 		Bucket: "node1",
// 		Org:    "miot-tracer",
// 	}
// )

// // log 配置
// var LogAddr = "/home/ianhui/code/miot-tracer/logFile/miot_tracer_log.json"

//-------------------以下是通用配置-------------------

var (
	//redis key前缀
	SecondIndex_prefix = "Second_index_"
	ThirdIndex_prefix  = "Third_index_"
	Node_prefix        = "node_"
	TAXI_ID_PREFIX     = "taxi_"
	//influxdb key前缀
	BucketNode_prefix = "node"
	//最大重试次数
	RETRY = 50
)

var (
	REF_TIME = time.Date(2008, 5, 17, 0, 0, 0, 0, time.UTC) // 我看sfs数据是从5月18号开始 ，那么这里定为5月17号
)

const (
	BIN_LEN         = int64(60 * 60 * 24) // 1 day
	BIN_BINARY_LEN  = 0x1f
	ELEMENTCODE_LEN = 0x7ff

	TS_LEN             = 16
	SEGMENT_LEN        = 16
	NEXT_NODE_LEN      = 8
	VARIABLE_CHECK_LEN = 15

	TYPE_SECOND_INDEX_FIRSTLINE SecondIndexType = "MIOT_SECONDE_INDEX_FIRSTLINE"
	TYPE_SECOND_INDEX_OTHERLINE SecondIndexType = "MIOT_SECONDE_INDEX_OTHERLINE"
)

const (
	THIRD_INDEX_NODE_ID_LEN      = 8
	THIRD_INDEX_SEQUENCE_NUM_LEN = 16
	TS_SKIP                      = 60 * 60 // 1h
	MAX_LEVEL                    = 64
)

// 一致性哈希配置
const (
	NODE_TOTAL_NUM = 16
)
