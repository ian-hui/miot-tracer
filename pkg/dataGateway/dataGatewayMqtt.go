package datagateway

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	iotlog.Infoln("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	iotlog.Errorf("Connect lost: %v", err)
}

type DataGatewayMqtt struct {
	mqttClient *mqtt.Client
}

func NewDataGatewayMqtt() *DataGatewayMqtt {
	opts := getMqttDefaultConfig()
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		//重试
		iotlog.Errorln("mqtt connect failed, retrying...")
		recoverd := false
		//如果连接失败，重试3次，如果还是失败，就退出
		for i := 0; i < RETRY_TIMES; i++ {
			if token := client.Connect(); token.Wait() && token.Error() != nil {
				iotlog.Errorln("mqtt connect failed, retrying...")
				time.Sleep(time.Second * 1)
			} else {
				recoverd = true
				break
			}
		}
		if !recoverd {
			iotlog.Fatalln("mqtt connect failed, exit. error: ", token.Error())
			panic("mqtt connect failed, exit...")
		}
	}
	return &DataGatewayMqtt{
		mqttClient: &client,
	}
}

func (d *DataGatewayMqtt) mqttPublish(topic string, message []byte) error {
	qos := QOS
	payload := message
	if token := (*d.mqttClient).Publish(topic, byte(qos), false, payload); token.Wait() && token.Error() != nil {
		iotlog.Errorln("publish failed")
		return token.Error()
	} else {
		iotlog.Infof("publish success, topic: %s, payload: %s\n", topic, payload)
	}
	return nil
}

func (d *DataGatewayMqtt) mqttSubscribe(topic string, handler func(client mqtt.Client, msg mqtt.Message)) {
	qos := QOS
	(*d.mqttClient).Subscribe(topic, byte(qos), handler)
}

func (d *DataGatewayMqtt) mqttDisconnect() {
	(*d.mqttClient).Disconnect(250)
}
