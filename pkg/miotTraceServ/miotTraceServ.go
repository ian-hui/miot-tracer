package miottraceserv

import (
	"encoding/json"
	mttypes "miot_tracing_go/mtTypes"
	"miot_tracing_go/pkg/logger"
	dataprocessor "miot_tracing_go/pkg/miotTraceServ/dataProcessor"
	indexprocessor "miot_tracing_go/pkg/miotTraceServ/indexProcessor"
)

var (
	iotlog = logger.Miotlogger
)



type MiotTracingServ interface {
}

type MiotTracingServImpl struct {
	dp *dataprocessor.DataProcessor
	ip *indexprocessor.IndexProcessor
	processChan chan mttypes.Message
}

func NewMiotTracingServImpl() *MiotTracingServImpl {
	return &MiotTracingServImpl{
		dp: dataprocessor.NewDataProcessor(),
		ip: indexprocessor.NewIndexProcessor(),
		processChan: make(chan mttypes.Message, 100),
	}
}



func (mts *MiotTracingServImpl) Start(worker_num int){
	for i := 0; i < worker_num; i++ {
		go mts.worker(i)
	}
}

func (m *MiotTracingServImpl) worker(workerID int) {
    for message := range m.processChan {
        // 处理消息
        m.handleMessage(message)
    }
}

func (m *MiotTracingServImpl) handleMessage(message mttypes.Message) {
	switch message.Type {
	case mttypes.TYPE_FIRST_UPLOAD:
		m.handleFirstData(message)
	case mttypes.TYPE_UPLOAD_INDEX:
		m.handleData(message)
	case mttypes.TYPE_UPLOAD_TAXI_DATA:
		m.handleIndex(message)
	default:
		iotlog.Errorf("unknown message type: %v", message.Type)
	}
}

func (m *MiotTracingServImpl) handleFirstData(message mttypes.Message) {
	var firstData mttypes.FirstData
	if err := json.Unmarshal(message.Content, &firstData); err != nil {
		iotlog.Errorf("json unmarshal error: %v", err)
		return
	}
	if firstData.TaxiFrontNode == "" {
		iotlog.Errorf("empty TaxiFrontNode")
		return
	}
	if firstData.Segment == 1{
		//add first data
		//增加索引
		ip := m.ip.
	}

}
