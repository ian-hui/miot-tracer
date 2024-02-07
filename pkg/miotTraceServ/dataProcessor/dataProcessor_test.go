package dataprocessor

import (
	"fmt"
	mttypes "miot_tracing_go/mtTypes"
	"testing"
	"time"
)

func TestInsertInfluxdb(t *testing.T) {
	dp := NewDataProcessor()
	defer dp.FlushData()
	d := &mttypes.TaxiData{
		TaxiID:    "1",
		Timestamp: "1213037726",
		Longitude: "3",
		Latitude:  "4",
		Occupancy: "5",
		TaxiDataLabel: mttypes.TaxiDataLabel{
			Segment: "1",
		},
	}
	err := dp.InsertTaxiData(d)
	if err != nil {
		t.Error(err)
	}
}

func TestQueryInfluxdb(t *testing.T) {
	dp := NewDataProcessor()
	// defer dp.flushData()
	// d := &mttypes.TaxiData{
	// 	TaxiID:    "5",
	// 	Timestamp: "1213537727",
	// 	Longitude: "3",
	// 	Latitude:  "4",
	// 	Occupancy: "5",
	// 	TaxiDataLabel: mttypes.TaxiDataLabel{
	// 		Segment: "2",
	// 	},
	// }
	// err := dp.insertTaxiData(d)
	// if err != nil {
	// 	t.Error(err)
	// }
	query := &mttypes.QueryStru{
		StartTime: "0",
		EndTime:   "2215037729",
		ID:        "5",
	}
	result, err := dp.QueryTaxiData(query)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("result:", result)
	for result.Next() {
		t.Log("record")
		t.Log(result.Record().Field())
		t.Log(result.Record().Values())
	}
}

func TestExt2Time(t *testing.T) {
	ext := int64(63349134527) // Unix时间戳
	wall := int64(0)          // 时间戳

	// 使用Unix时间戳（ext）创建时间
	extTime := time.Unix(ext, 0)

	// 使用时间戳（wall）创建时间
	wallTime := time.Unix(0, wall)

	fmt.Println("Unix时间戳转换为时间:", extTime)
	fmt.Println("时间戳转换为时间:", wallTime)
}
