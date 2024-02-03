package mttypes

type TaxiData struct {
	TaxiID    string
	Timestamp string
	Longitude string
	Latitude  string
	Occupancy string
	TaxiDataLabel
}

type TaxiDataLabel struct {
	Segment int
}

type TaxiFrontNode string

type TaxiInfo struct {
	TaxiID      string
	Index       []string
	FronterNode string
	Segment     string
}

type FirstData struct {
	TaxiData
	TaxiFrontNode
}

type Message struct {
	// 消息对应事件的类型, 详情见《事件类型》
	Type string `json:"type"`
	// 消息内容
	Content []byte `json:"content"`
}

// 事件类型
var (
	TYPE_FIRST_UPLOAD     = "MIOT_FIRST_UPLOAD"     //第一次上传数据
	TYPE_UPLOAD_TAXI_DATA = "MIOT_UPLOAD_TAXI_DATA" //上传出租车数据
	TYPE_UPLOAD_INDEX     = "MIOT_UPLOAD_INDEX"     //用于跳表索引
)

// metadata
// ID ：starttime，endtime，segment，nextNode
// 每个id做一个list
type SecondIndex struct {
	ID            string `json:"id"`
	StartTs       string `json:"startts"`
	EndTs         string `json:"endts"`
	Segment       string `json:"segment"`
	NextNode      string `json:"nextnode"`
	NextMetaIndex string `json:"nextmetaindex"` //下一个metadata的在redislist的索引
}

// index
// id : [timestamp,nodeid,segment]
type ThirdIndex struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	NodeID    string `json:"nodeid"`
	Segment   string `json:"segment"`
}

type RedisConf struct {
	Addr    string `json:"addr"`
	Pwd     string `json:"pwd"`
	DB      string `json:"db"`
	Timeout string `json:"timeout"`
}

type InfluxConf struct {
	Url    string `json:"url"`
	Token  string `json:"token"`
	Bucket string `json:"bucket"`
	Org    string `json:"org"`
}

type QueryStru struct {
	StartTime string `json:"starttime"`
	EndTime   string `json:"endtime"`
	Segment   string `json:"segment"`
	ID        string `json:"id"`
}

type SecondIndexType string
