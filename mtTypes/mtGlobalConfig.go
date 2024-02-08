package mttypes

import (
	"fmt"
	"os"
	"time"
)

func init() {
	fmt.Println("init global config")
	os.Setenv("NODE_ID", "1")
	os.Setenv("INFLUXDB_TOKEN", "J_xeoyLkPQFHBilXk4ELHjV85A7fFtIJvlo3GTjmKnF3QPZU63H7N0FH5_x7JBMPy3MRvVwoeoW0rnReDyLuPg==")
	os.Setenv("INFLUXDB_URL", "http://localhost:8086")
	os.Setenv("INFLUXDB_BUCKET", "node1")
	os.Setenv("REDIS_URL", "localhost:6379")
}

// var NODE_ID = os.Getenv("NODE_ID")
var NODE_ID string = "1"

// redis配置
var (
	RedisConfig = RedisConf{
		Addr:    "localhost:6379",
		Pwd:     "reins5401",
		DB:      "0",
		Timeout: "10",
	}
	SecondIndex_prefix = "Second_index_"
	ThirdIndex_prefix  = "Third_index_"
	Node_prefix        = "node_"
	TAXI_ID_PREFIX     = "taxi_"
)

// influxdb配置
var (
	InfluxConfig = InfluxConf{
		Url:    "http://localhost:8086",
		Token:  "J_xeoyLkPQFHBilXk4ELHjV85A7fFtIJvlo3GTjmKnF3QPZU63H7N0FH5_x7JBMPy3MRvVwoeoW0rnReDyLuPg==",
		Bucket: "node1",
		Org:    "miot-tracer",
	}
	BucketNode_prefix = "node"
)

var (
	REF_TIME = time.Date(2008, 1, 1, 0, 0, 0, 0, time.UTC)
)

const (
	BIN_LEN            = int64(60 * 60 * 24) // 1 day
	BIN_BINARY_LEN     = 0x1f
	ELEMENTCODE_LEN    = 0x7ff
	TS_LEN             = 0xffff
	SEGMENT_LEN        = 0xff
	NEXT_NODE_LEN      = 0xff
	VARIABLE_CHECK_LEN = (1 << 15) - 1

	TYPE_SECOND_INDEX_FIRSTLINE SecondIndexType = "MIOT_SECONDE_INDEX_FIRSTLINE"
	TYPE_SECOND_INDEX_OTHERLINE SecondIndexType = "MIOT_SECONDE_INDEX_OTHERLINE"
)

const (
	THIRD_INDEX_NODE_ID_LEN      = 8
	THIRD_INDEX_SEQUENCE_NUM_LEN = 16
	TS_SKIP                      = 60 * 30 // 30 min
)

// 一致性哈希配置
const (
	NODE_TOTAL_NUM = 3
)
