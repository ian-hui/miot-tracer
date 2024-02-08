package datagateway

import (
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
	processChan chan mttypes.Message
	sendingChan chan mttypes.Message
}

func NewDataGateway(m miottraceserv.MiotTracingServ) Datagateway {

	return &dataGatewayImpl{
		mqtt:        NewDataGatewayMqtt(),
		http:        NewDataGatewayHTTP(),
		processChan: m.GetProcessChan(),
		sendingChan: m.GetSendingChan(),
	}
}

func (d dataGatewayImpl) Start() {
	// mqtt
	d.mqtt.mqttSubscribe(RECEIVER, d.RECEIVERHandler)
	go d.asyncSendMsg() //异步发送消息
	// http
	d.http.Start()
}
