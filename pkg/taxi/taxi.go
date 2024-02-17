package taxi

import (
	"fmt"
	"strconv"
)

type Taxi interface {
	CheckRegion(longitude, latitude string) (string, error)
	NearestPointInRegion(lon, lat string) (float64, float64)
}

type taxi struct {
	northest string
	southest string
	westest  string
	eastest  string
}

func NewTaxi() Taxi {

	// 三藩市中心的经纬度范围
	return &taxi{
		northest: "37.8379",
		southest: "37.7379",
		eastest:  "-122.3575",
		westest:  "-122.4575"}
}

func (dp *taxi) CheckRegion(longitude, latitude string) (string, error) {
	// Parse the input strings to float64
	long, err1 := strconv.ParseFloat(longitude, 64)
	lat, err2 := strconv.ParseFloat(latitude, 64)
	if err1 != nil || err2 != nil {
		return "invalid input", err1
	}

	westest, err := strconv.ParseFloat(dp.westest, 64)
	if err != nil {
		return "invalid input", err
	}
	eastest, err := strconv.ParseFloat(dp.eastest, 64)
	if err != nil {
		return "invalid input", err
	}
	southest, err := strconv.ParseFloat(dp.southest, 64)
	if err != nil {
		return "invalid input", err
	}
	northest, err := strconv.ParseFloat(dp.northest, 64)
	if err != nil {
		return "invalid input", err
	}
	// 判断是否在市区
	if long < westest || long > eastest || lat < southest || lat > northest {
		return "out", nil
	}

	// 计算每个正方形的大小
	longStep := (eastest - westest) / 4
	latStep := (northest - southest) / 4

	// 计算区域ID
	longIndex := int((long - westest) / longStep)
	latIndex := int((lat - southest) / latStep)

	// 生成区域ID（从左到右，从上到下）
	regionID := latIndex*4 + longIndex + 1

	return fmt.Sprintf("%d", regionID), nil
}

func (dp *taxi) NearestPointInRegion(longitude, latitude string) (float64, float64) {
	var nearest_lon, nearest_lat float64
	// 计算每个正方形的大小
	west, _ := strconv.ParseFloat(dp.westest, 64)
	east, _ := strconv.ParseFloat(dp.eastest, 64)
	south, _ := strconv.ParseFloat(dp.southest, 64)
	north, _ := strconv.ParseFloat(dp.northest, 64)
	lon, _ := strconv.ParseFloat(longitude, 64)
	lat, _ := strconv.ParseFloat(latitude, 64)
	if lon < west {
		nearest_lon = west
	} else if lon > east {
		nearest_lon = east
	} else {
		nearest_lon = lon
	}

	if lat < south {
		nearest_lat = south
	} else if lat > north {
		nearest_lat = north
	} else {
		nearest_lat = lat
	}

	return nearest_lon, nearest_lat
}
