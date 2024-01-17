package datagateway

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type dataGatewayHTTP struct {
	gin *gin.Engine
}

func (d *dataGatewayHTTP) Start() {
	err := d.gin.Run(":8080")
	if err != nil {
		//重试
		iotlog.Errorln("http server start failed, retrying...")
		recoverd := false
		//如果连接失败，重试3次，如果还是失败，就退出
		for i := 0; i < RETRY_TIMES; i++ {
			err := d.gin.Run(":8080")
			if err != nil {
				iotlog.Errorln("http server start failed, retrying...")
			} else {
				recoverd = true
				break
			}
		}
		if !recoverd {
			iotlog.Fatalln("http server start failed, exit...")
			panic("http server start failed, exit...")
		}
	}
}

//-------------------- controller&router --------------------
//-------------------- controller&router --------------------
//-------------------- controller&router --------------------

func NewDataGatewayHTTP() *dataGatewayHTTP {
	router := gin.Default()
	router.GET("/search/:taxi_id/:start_time/:end_time", searchHandler)
	return &dataGatewayHTTP{
		gin: router,
	}
}

func searchHandler(c *gin.Context) {
	taxiID := c.Param("taxi_id")
	startTime := c.Param("start_time")
	endTime := c.Param("end_time")
	//查询
	fmt.Println(taxiID, startTime, endTime)
	data := "taxi_id,timestamp,longitude,latitude\n"
	//返回
	returnCSV(c, data)
}

//-------------------- return格式 --------------------
//-------------------- return格式 --------------------
//-------------------- return格式 --------------------

func returnCSV(c *gin.Context, data string) {
	// 设置响应头
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=myfile.csv")
	c.Header("Content-Type", "application/csv")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// 发送 CSV 数据
	c.String(http.StatusOK, data)
}

func returnError(c *gin.Context, err error, errcode int) {
	c.JSON(errcode, gin.H{
		"error": err.Error(),
	})
}
