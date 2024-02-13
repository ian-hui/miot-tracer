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
	// d := &mttypes.TaxiData{
	// 	TaxiID:    "1",
	// 	Timestamp: "1213537727",
	// 	Longitude: "3",
	// 	Latitude:  "4",
	// 	Occupancy: "5",
	// 	TaxiDataLabel: mttypes.TaxiDataLabel{
	// 		Segment: "2",
	// 	},
	// }
	// err := dp.InsertTaxiData(d)
	// if err != nil {
	// 	t.Error(err)
	// }
	// dp.FlushData()

	// time := time.Now().UTC().Unix()
	// time_string := strconv.FormatInt(time, 10)
	query := &mttypes.QueryStru{
		StartTime: "0",
		EndTime:   "1213537790",
		ID:        "1",
		Segment:   "2",
	}
	result, err := dp.QueryTaxiData(query)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("result:", result)

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
