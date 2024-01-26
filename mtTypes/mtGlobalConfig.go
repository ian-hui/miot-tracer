package mttypes

// var NODE_ID = os.Getenv("NODE_ID")
var NODE_ID = "1"

// redis配置
var (
	RedisConfig = RedisConf{
		Addr:    "localhost:6379",
		Pwd:     "reins5401",
		DB:      "0",
		Timeout: "10",
	}
	Meta_prefix  = "meta_"
	Index_prefix = "index_"
	Node_prefix  = "node_"
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
