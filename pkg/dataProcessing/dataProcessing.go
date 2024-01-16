package dataprocessing

import (
	"fmt"
	"strconv"
)

type DataProcessing interface {
	CheckRegion(longitude, latitude string) (string, error)
}

type dataProcessing struct {
	northest string
	southest string
	westest  string
	eastest  string
}

func NewDataProcessing() DataProcessing {

	// 三藩市中心的经纬度范围
	return &dataProcessing{
		northest: "37.8379",
		southest: "37.7379",
		eastest:  "-122.4575",
		westest:  "-122.3575"}
}

func (dp *dataProcessing) CheckRegion(longitude, latitude string) (string, error) {
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
	if !(long < westest && long > eastest && lat > southest && lat < northest) {
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
