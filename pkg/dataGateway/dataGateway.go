package datagateway

import (
	"miot_tracing_go/pkg/logger"
)

var (
	iotlog = logger.Miotlogger
)

type Datagateway interface {
	Start()
}
type dataGatewayImpl struct {
}

func NewDataGateway() Datagateway {
	return &dataGatewayImpl{}
}

func (d dataGatewayImpl) Start() {
	// mqtt
	dataGatewayMqtt := NewDataGatewayMqtt()
	dataGatewayMqtt.mqttSubscribe(TopicUploadIndex, uploadIndexHandler)
	dataGatewayMqtt.mqttSubscribe(TopicUploadTaxiData, uploadTaxiDataHandler)
	// http
	dataGatewayHttp := NewDataGatewayHTTP()
	dataGatewayHttp.Start()
}
