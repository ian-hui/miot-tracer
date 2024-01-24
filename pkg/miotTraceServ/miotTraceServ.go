package miottraceserv

import "miot_tracing_go/pkg/logger"

var (
	iotlog = logger.Miotlogger
)

//如果segment没见过且不是1，
//那么回传给前一个节点
//在本节点创建一个索引
