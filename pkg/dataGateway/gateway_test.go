package datagateway

import (
	"fmt"
	"testing"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func TestMqtt(T *testing.T) {
	c := NewDataGatewayMqtt()
	defer c.mqttDisconnect()
	c.mqttSubscribe("test", func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	})
	c.mqttPublish("test", []byte("hello world"))
	select {}
}

func TestHttp(T *testing.T) {
	c := NewDataGatewayHTTP()
	c.Start()
}
