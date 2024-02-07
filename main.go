package main

import (
	datagateway "miot_tracing_go/pkg/dataGateway"
	miottraceserv "miot_tracing_go/pkg/miotTraceServ"
)

func main() {
	// time.Sleep(10 * time.Hour)
	mts := miottraceserv.NewMiotTracingServImpl()
	gw := datagateway.NewDataGateway(mts)
	mts.Start(4)
	gw.Start()
}
