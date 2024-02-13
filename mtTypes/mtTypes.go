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
	TaxiID      string   `json:"id"`
	Index       []string `json:"index"`
	FronterNode string   `json:"fronter_node"`
	Segment     string   `json:"segment"`
	PreTime     string   `json:"pre_time"` // 由于希望taxi半小时发送一次信息，所以🚕自身需要记录一个上次发送信息的时间，一旦超过半小时，就发送一次upload_third_index
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
	Content interface{} `json:"content"`
}

// 事件类型
var (
	TYPE_FIRST_UPLOAD            = "MIOT_FIRST_UPLOAD"            //第一次上传数据
	TYPE_UPLOAD_TAXI_DATA        = "MIOT_UPLOAD_TAXI_DATA"        //上传出租车数据
	TYPE_UPLOAD_THIRD_INDEX      = "MIOT_UPLOAD_THIRD_INDEX"      //接收🚕的索引并转发
	TYPE_UPLOAD_THIRD_INDEX_HEAD = "MIOT_UPLOAD_THIRD_INDEX_HEAD" //存储头节点
	TYPE_UPDATE_SECOND_INDEX     = "MIOT_UPDATE_SECOND_INDEX"     //更新二级索引(补全endtime和nextnode)
	TYPE_UPDATE_THIRD_INDEX      = "MIOT_UPDATE_THIRD_INDEX"      //更新三级索引(接收转发并存储)
	TYPE_BUILD_QUERY             = "MIOT_BUILD_QUERY"             //构建查询
	TYPE_SEARCH_THIRD_INDEX      = "MIOT_SEARCH_THIRD_INDEX"      //查询三级索引
	TYPE_QUERY_TAXI_DATA         = "MIOT_QUERY_TAXI_DATA"         //查询出租车数据
	TYPE_SEND_BACK_RESULT        = "MIOT_SEND_BACK_RESULT"        //返回查询结果
)

// metadata
// ID ：starttime，endtime，segment，nextNode
// 每个id做一个list
type SecondIndex struct {
	ID       string `json:"id"`
	StartTs  string `json:"startts"`
	EndTs    string `json:"endts"`
	Segment  string `json:"segment"`
	NextNode string `json:"nextnode"`
	// NextSecondIndex string `json:"nextnodeindex"` //下一个节点的二级索引在redis—list的index
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
	Tii         ThirdIndexInfo //查询条件
	TraverseCfg TraverseConfig //遍历有关的参数
	QueryNode   string         `json:"querynode"` //记录发起查询的节点，用于回传信息
	RequestID   string         `json:"requestid"` //记录用户终端id，用于回传信息
	StartTime   string         `json:"starttime"`
	EndTime     string         `json:"endtime"`
	Segment     string         `json:"segment"`
	ID          string         `json:"id"`
}

type ThirdIndexInfo struct {
	Taxi_Start_Ts string `json:"taxi_start_ts"`
	Skip_Ts       string `json:"skip_ts"`
}

type TraverseConfig struct {
	Start_segment    string   `json:"start_segment"`    //开始节点
	Traversed        []string `json:"traversed"`        //遍历过的节点
	Previous_segment string   `json:"previous_segment"` //上一个节点
}

type SecondIndexType string

type Result struct {
	Request_id      string                   `json:"request_id"`
	Start_segment   string                   `json:"start_segment"`
	Current_segment string                   `json:"current_segment"`
	End_segment     string                   `json:"end_segment"`
	Result          []map[string]interface{} `json:"result"`
}

// type SegmentType struct {
// 	Seg_type string `json:"seg_type"`
// 	Seg_id   string `json:"seg_id"`
// }

// const (
// 	SEGMENT_START   = "start"
// 	SEGMENT_CURRENT = "current"
// 	SEGMENT_END     = "end"
// )
