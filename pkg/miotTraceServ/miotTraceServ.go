package miottraceserv

import (
	"miot_tracing_go/pkg/logger"
	dataprocessor "miot_tracing_go/pkg/miotTraceServ/dataProcessor"
	indexprocessor "miot_tracing_go/pkg/miotTraceServ/indexProcessor"
)

var (
	iotlog = logger.Miotlogger
)

//如果segment没见过且不是1，
//那么回传给前一个节点
//在本节点创建一个索引

type MiotTracingServ interface {
}

type MiotTracingServImpl struct {
	dp *dataprocessor.DataProcessor
	ip *indexprocessor.IndexProcessor
}
