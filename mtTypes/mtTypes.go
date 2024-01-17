package mttypes

type TaxiData struct {
	TaxiID    string
	Timestamp string
	Longitude string
	Latitude  string
	Occupancy string
}

type TaxiInfo struct {
	TaxiID      string
	Index       []string
	FronterNode string
	Segment     []string
}

type Message struct {
	// 消息对应事件的类型, 详情见《事件类型》
	Type string `json:"type"`
	// 消息内容
	Content string `json:"content"`
}

// 事件类型
var (
	UPLOAD_TAXI_DATA = "MIOT_UPLOAD_TAXI_DATA"
	UPLOAD_INDEX     = "MIOT_UPLOAD_INDEX"
)
