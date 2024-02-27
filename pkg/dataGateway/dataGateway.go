package datagateway

import (
	"encoding/json"
	mttypes "miot_tracing_go/mtTypes"
	"miot_tracing_go/pkg/logger"
	miottraceserv "miot_tracing_go/pkg/miotTraceServ"
)

var (
	iotlog = logger.Miotlogger
)

type Datagateway interface {
	Start()
}
type dataGatewayImpl struct {
	mqtt        *DataGatewayMqtt
	http        *DataGatewayHTTP
	reuqest_map map[string]chan interface{}
	processChan chan mttypes.Message
	sendingChan chan mttypes.Message
	resultChan  chan mttypes.Result
}

func NewDataGateway(m miottraceserv.MiotTracingServ) Datagateway {

	return &dataGatewayImpl{
		mqtt:        NewDataGatewayMqtt(),
		http:        NewDataGatewayHTTP(),
		processChan: m.GetProcessChan(),
		sendingChan: m.GetSendingChan(),
		resultChan:  m.GetResultChan(),
	}
}

func (d dataGatewayImpl) Start() {
	// mqtt
	d.mqtt.mqttSubscribe(RECEIVER, d.RECEIVERHandler)
	go d.asyncSendMsg() //异步发送消息
	// go d.asyncResDeliver() //异步发送res
	// http
	d.http.gin.GET("/search/:taxi_id/:start_time/:end_time", d.searchHandler)
	d.http.Start()
}

func (d *dataGatewayImpl) asyncSendMsg() {
	for msg := range d.sendingChan {
		b_msg, err := json.Marshal(msg)
		if err != nil {
			iotlog.Errorln("marshal failed :", err)
			continue
		}
		err = d.mqtt.mqttPublish(msg.Topic, b_msg)
		if err != nil {
			iotlog.Errorln("publish failed :", err)
		}
	}
}

//接收到信息后放入一个map中，等待合并
// func (d *dataGatewayImpl) asyncResReceiver() {
// 	for result := range d.resultChan {
// 		request_id := result.Request_id
// 		d.reuqest_map[request_id] <- result
// 	}
// }

// func (d *dataGatewayImpl) asyncResDeliver() {
// 	for result := range d.resultChan {

// 	}
// }
