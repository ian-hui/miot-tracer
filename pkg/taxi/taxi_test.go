package taxi_test

import (
	"encoding/csv"
	"fmt"
	mttypes "miot_tracing_go/mtTypes"
	"miot_tracing_go/pkg/taxi"

	"os"
	"strconv"
	"testing"
)

var (
	increase = 0
)

func TestCheckRegion(t *testing.T) {
	filepath := "/home/ianhui/code/miot-tracer/datas/indexData/combined_taxi_data_sfs.csv"
	dp := taxi.NewTaxi()
	// 读取CSV文件
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	// 读取CSV文件标题行
	_, err = reader.Read()
	if err != nil {
		panic(err)
	}
	// 逐行读取
	for i := 0; i < 3000; i++ {
		record, err := reader.Read()
		if err != nil {
			break
		}

		// 将读取的行转换为TaxiData结构
		taxiData := mttypes.TaxiData{
			TaxiID:    record[0],
			Latitude:  record[1],
			Longitude: record[2],
			Occupancy: record[3],
			Timestamp: record[4],
		}

		// 检查是否在指定区域内
		region, err := dp.CheckRegion(taxiData.Longitude, taxiData.Latitude)
		if err != nil {
			fmt.Println("Error checking region:", err)
			return
		}

		if region == "out" {
			fmt.Println("before lon", taxiData.Longitude, "before lat", taxiData.Latitude, "region", region)
			//找寻区域内最近的点
			lon, lat := dp.NearestPointInRegion(taxiData.Longitude, taxiData.Latitude)
			taxiData.Longitude, taxiData.Latitude = strconv.FormatFloat(lon, 'f', -1, 64), strconv.FormatFloat(lat, 'f', -1, 64)
			region, err = dp.CheckRegion(taxiData.Longitude, taxiData.Latitude)
			if err != nil {
				fmt.Println("Error checking region:", err)
				return
			}
			fmt.Println("after lon", taxiData.Longitude, "after lat", taxiData.Latitude, "region", region)
		}
		id, err := strconv.Atoi(taxiData.TaxiID)
		if err != nil {
			fmt.Println("Error checking region:", err)
			return
		}
		id = id + increase
		fmt.Println("timestamp", taxiData.Timestamp, "taxiid", id, "Longitude",
			taxiData.Longitude, "latitude", taxiData.Latitude, "region", region)
	}
}
