package datagateway

import (
	"encoding/json"
	mttypes "miot_tracing_go/mtTypes"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (d *dataGatewayImpl) RECEIVERHandler(client mqtt.Client, msg mqtt.Message) {
	var message mttypes.Message
	err := json.Unmarshal(msg.Payload(), &message)
	if err != nil {
		iotlog.Errorln("unmarshal failed :", err)
	}
	d.processChan <- message
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
