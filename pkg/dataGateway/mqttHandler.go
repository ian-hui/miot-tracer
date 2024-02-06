package datagateway

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func RECEIVERHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received `%s` from `%s` topic", msg.Payload(), msg.Topic())
}
