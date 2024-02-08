package main

import (
	mttypes "miot_tracing_go/mtTypes"
	datagateway "miot_tracing_go/pkg/dataGateway"
	miottraceserv "miot_tracing_go/pkg/miotTraceServ"
	"testing"
)

func TestServHandleFirst(t *testing.T) {
	mts := miottraceserv.NewMiotTracingServImpl()
	gw := datagateway.NewDataGateway(mts)

	go func() {
		first := mttypes.FirstData{
			TaxiData: mttypes.TaxiData{
				TaxiID: "1",
				//unix时间戳
				Timestamp: "1199449857",
				//经度
				Longitude: "123.123",
				//纬度
				Latitude:  "123.123",
				Occupancy: "1",
				TaxiDataLabel: mttypes.TaxiDataLabel{
					Segment: "2",
				},
			},
			TaxiFrontNode: "3",
		}

		message := mttypes.Message{
			Type:    mttypes.TYPE_FIRST_UPLOAD,
			Content: first,
		}

		mts.HandleFirstData(message)
	}()

	mts.Start(4)
	gw.Start()
}
