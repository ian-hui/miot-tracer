package dataprocessor

import (
	"context"
	"fmt"
	mttypes "miot_tracing_go/mtTypes"
	"miot_tracing_go/pkg/logger"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

var (
	client influxdb2.Client
	iotlog = logger.Miotlogger
)

type DataProcessor struct {
	dbclient influxdb2.Client
}

func NewDataProcessor() *DataProcessor {
	if client == nil {
		c := influxdb2.NewClient(mttypes.InfluxConfig.Url, mttypes.InfluxConfig.Token)
		client = c
	}
	return &DataProcessor{client}
}

func (dp *DataProcessor) InsertTaxiData(d *mttypes.TaxiData) error {
	writeAPI := client.WriteAPI(mttypes.InfluxConfig.Org, mttypes.InfluxConfig.Bucket)
	ts, err := CSVts2timestamp(d.Timestamp)
	if err != nil {
		iotlog.Errorln("CSVts2timestamp failed, err:", err)
		return err
	}
	p := influxdb2.NewPoint("taxi",
		map[string]string{"taxi_id": d.TaxiID, "segment": d.Segment},
		map[string]interface{}{"longitude": d.Longitude, "latitude": d.Latitude, "occupancy": d.Occupancy},
		ts)
	fmt.Println("p:", p)
	writeAPI.WritePoint(p)
	return nil
}

func (dp *DataProcessor) FlushData() {
	writeAPI := client.WriteAPI(mttypes.InfluxConfig.Org, mttypes.InfluxConfig.Bucket)
	writeAPI.Flush()
}

func (dp *DataProcessor) QueryTaxiData(query *mttypes.QueryStru) (result *api.QueryTableResult, err error) {
	queryAPI := client.QueryAPI(mttypes.InfluxConfig.Org)
	// 准备参数
	bucketName := mttypes.BucketNode_prefix + mttypes.NODE_ID
	StartTime, err := unixTimestamp_2_RFC3339(query.StartTime)
	if err != nil {
		iotlog.Errorln("unixTimestamp_2_RFC3339 failed, err:", err)
		return
	}
	EndTime, err := unixTimestamp_2_RFC3339(query.EndTime)
	if err != nil {
		iotlog.Errorln("unixTimestamp_2_RFC3339 failed, err:", err)
		return
	}
	// 查询 |> group(columns: ["segment"])
	result, err = queryAPI.Query(context.Background(),
		`from(bucket:"`+bucketName+`")
			|> range(start:`+StartTime+`,stop:`+EndTime+`) 
			|> filter(fn: (r) => r._measurement == "taxi" )
			|> filter(fn: (r) => r["taxi_id"] == "`+query.ID+`")
			|> filter(fn: (r) => r["segment"] == "`+query.Segment+`")
	`)
	if err != nil {
		iotlog.Errorln("query failed, err:", err)
		return
	}

	return
}

func (dp *DataProcessor) ClientClose() {
	client.Close()
}

// --------------------helper functions--------------------
func CSVts2timestamp(ts string) (time.Time, error) {
	timestamp, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(timestamp, 0).UTC(), nil
}

func unixTimestamp_2_RFC3339(ts string) (string, error) {
	i, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		iotlog.Errorln("strconv.ParseInt failed, err:", err)
		return "", err
	}
	utcTime := time.Unix(i, 0).UTC()
	return utcTime.Format(time.RFC3339), nil
}
