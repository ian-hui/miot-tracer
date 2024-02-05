package mttypes

import "time"

// var NODE_ID = os.Getenv("NODE_ID")
const NODE_ID string = "1"

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
)

// influxdb配置
var (
	InfluxConfig = InfluxConf{
		Url:    "http://localhost:8086",
		Token:  "1pJyO11iMPsN4MQ0E-gkVcq5ZlgSyvsjiiMgUBGn9rmGfYr3TU3ekMNpJTwbpK15dkVDelPS6nYeM7eRwWBSVg==",
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
