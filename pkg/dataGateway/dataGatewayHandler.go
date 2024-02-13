package datagateway

import (
	"encoding/json"
	mttypes "miot_tracing_go/mtTypes"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (d *dataGatewayImpl) RECEIVERHandler(client mqtt.Client, msg mqtt.Message) {
	var message mttypes.Message
	err := json.Unmarshal(msg.Payload(), &message)
	if err != nil {
		iotlog.Errorln("unmarshal failed :", err)
	}
	d.processChan <- message
}

func (d *dataGatewayImpl) searchHandler(c *gin.Context) {
	taxiID := c.Param("taxi_id")
	startTime := c.Param("start_time")
	endTime := c.Param("end_time")
	request_id := uuid.New().String()
	//注册
	d.reuqest_map[request_id] = make(chan interface{})
	//发送
	d.sendingChan <- mttypes.Message{
		Type: mttypes.TYPE_QUERY_TAXI_DATA,
		//build mttypes.QueryStru
		Content: mttypes.QueryStru{
			ID:        taxiID,
			StartTime: startTime,
			EndTime:   endTime,
			RequestID: request_id,
			QueryNode: mttypes.NODE_ID,
		},
	}
	//等待
	result := <-d.reuqest_map[request_id]
	//返回
	returnCSV(c, result.(string))
}
