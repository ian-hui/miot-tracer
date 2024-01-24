package mttypes

// var NODE_ID = os.Getenv("NODE_ID")
var NODE_ID = "1"
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
