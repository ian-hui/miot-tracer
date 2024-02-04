package datagateway

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func uploadTaxiDataHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received `%s` from `%s` topic", msg.Payload(), msg.Topic())
}

func upload3IndexHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received `%s` from `%s` topic", msg.Payload(), msg.Topic())
}

func update2IndexHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received `%s` from `%s` topic", msg.Payload(), msg.Topic())
}
