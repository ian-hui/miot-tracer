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
	processChan chan mttypes.Message
	sendingChan chan mttypes.Message
}

func NewDataGateway(miottraceserv.MiotTracingServ) Datagateway {
	return &dataGatewayImpl{
		processChan: miottraceserv.NewMiotTracingServImpl().GetProcessChan(),
		sendingChan: miottraceserv.NewMiotTracingServImpl().GetSendingChan(),
	}
}

func (d dataGatewayImpl) Start() {
	// mqtt
	dataGatewayMqtt := NewDataGatewayMqtt()
	dataGatewayMqtt.mqttSubscribe(TopicUpload3Index, upload3IndexHandler)
	dataGatewayMqtt.mqttSubscribe(TopicUploadTaxiData, uploadTaxiDataHandler)
	dataGatewayMqtt.mqttSubscribe(TopicUpdate2Index, update2IndexHandler)
	// http
	dataGatewayHttp := NewDataGatewayHTTP()
	dataGatewayHttp.Start()
}
