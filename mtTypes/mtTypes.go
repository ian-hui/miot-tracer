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
	Segment string
}

type TaxiFrontNode string

type TaxiInfo struct {
	TaxiID      string
	Index       []string
	FronterNode string
	Segment     string
	PreTime     string
}

type FirstData struct {
	TaxiData
	TaxiFrontNode
}

type Message struct {
	// 消息对应事件的类型, 详情见《事件类型》
	Type  string `json:"type"`
	Topic string `json:"topic"`
	// 消息内容
	Content []byte `json:"content"`
}

// 事件类型
var (
	TYPE_FIRST_UPLOAD            = "MIOT_FIRST_UPLOAD"            //第一次上传数据
	TYPE_UPLOAD_TAXI_DATA        = "MIOT_UPLOAD_TAXI_DATA"        //上传出租车数据
	TYPE_UPLOAD_THIRD_INDEX      = "MIOT_UPLOAD_THIRD_INDEX"      //接收🚕的索引并转发
	TYPE_UPLOAD_THIRD_INDEX_HEAD = "MIOT_UPLOAD_THIRD_INDEX_HEAD" //存储头节点
	TYPE_UPDATE_SECOND_INDEX     = "MIOT_UPDATE_SECOND_INDEX"     //更新二级索引(补全endtime和nextnode)
	TYPE_UPDATE_THIRD_INDEX      = "MIOT_UPDATE_THIRD_INDEX"      //更新三级索引(接收转发并存储)
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
	ID          string `json:"id"`
	SequenceNum string `json:"sequencenum"`
	NodeID      string `json:"nodeid"`
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
	Tii       ThirdIndexInfo //查询条件
	StartTime string         `json:"starttime"`
	EndTime   string         `json:"endtime"`
	Segment   string         `json:"segment"`
	ID        string         `json:"id"`
}

type ThirdIndexInfo struct {
	Taxi_Start_Ts string `json:"taxi_start_ts"`
	Skip_Ts       string `json:"skip_ts"`
}

type SecondIndexType string
