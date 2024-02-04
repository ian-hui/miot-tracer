package miottraceserv

import (
	"encoding/json"
	mttypes "miot_tracing_go/mtTypes"
	"miot_tracing_go/pkg/logger"
	dataprocessor "miot_tracing_go/pkg/miotTraceServ/dataProcessor"
	indexprocessor "miot_tracing_go/pkg/miotTraceServ/indexProcessor"
	"strconv"
)

var (
	iotlog = logger.Miotlogger
)

type MiotTracingServ interface {
	Start(worker_num int)
	GetProcessChan() chan mttypes.Message
	GetSendingChan() chan mttypes.Message
}

type MiotTracingServImpl struct {
	dp          *dataprocessor.DataProcessor
	ip          *indexprocessor.IndexProcessor
	processChan chan mttypes.Message
	sendingChan chan mttypes.Message
}

func NewMiotTracingServImpl() *MiotTracingServImpl {
	return &MiotTracingServImpl{
		dp:          dataprocessor.NewDataProcessor(),
		ip:          indexprocessor.NewIndexProcessor(),
		processChan: make(chan mttypes.Message, 100),
		sendingChan: make(chan mttypes.Message, 100),
	}
}

func (mts *MiotTracingServImpl) Start(worker_num int) {
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

func (m *MiotTracingServImpl) handleMessage(message mttypes.Message) error {
	switch message.Type {
	case mttypes.TYPE_FIRST_UPLOAD:
		m.handleFirstData(message)
	case mttypes.TYPE_UPLOAD_THIRD_INDEX:
		m.handleThirdIndex(message)
	case mttypes.TYPE_UPLOAD_TAXI_DATA:
		m.handleData(message)
	case mttypes.TYPE_UPDATE_SECOND_INDEX:
		m.handleUploadSecondIndex(message)
	default:
		iotlog.Errorf("unknown message type: %v", message.Type)
	}
}

func (m *MiotTracingServImpl) handleFirstData(message mttypes.Message) (err error) {
	var firstData mttypes.FirstData
	if err = json.Unmarshal(message.Content, &firstData); err != nil {
		iotlog.Errorf("json unmarshal error: %v", err)
		return
	}
	if firstData.TaxiFrontNode == "" {
		iotlog.Errorf("empty TaxiFrontNode")
		return
	}
	//add first data
	//增加索引
	err = m.ip.SIP.CreateSecondIndex(&mttypes.SecondIndex{
		ID:      firstData.TaxiID,
		StartTs: firstData.Timestamp,
		Segment: firstData.Segment,
	})
	if err != nil {
		iotlog.Errorf("CreateSecondIndex error: %v", err)
		return
	}

	//存储数据
	err = m.dp.InsertTaxiData(&mttypes.TaxiData{
		TaxiID:    firstData.TaxiID,
		Timestamp: firstData.Timestamp,
		Longitude: firstData.Longitude,
		Latitude:  firstData.Latitude,
		Occupancy: firstData.Occupancy,
		TaxiDataLabel: mttypes.TaxiDataLabel{
			Segment: firstData.Segment,
		},
	})
	if err != nil {
		iotlog.Errorf("InsertTaxiData error: %v", err)
		return
	}
	if firstData.Segment != "1" {
		//把endtime传输给channel，传输到前一个节点
		sendback_second_index := mttypes.SecondIndex{
			ID:       firstData.TaxiID,
			EndTs:    firstData.Timestamp,
			Segment:  decreaseSegment(firstData.Segment),
			NextNode: mttypes.NODE_ID,
		}
		sendback_second_index_json, err := json.Marshal(sendback_second_index)
		if err != nil {
			iotlog.Errorf("json marshal error: %v", err)
			return err
		}
		m.sendingChan <- mttypes.Message{
			Type:    mttypes.TYPE_UPDATE_SECOND_INDEX,
			Content: sendback_second_index_json,
		}
	}
	return
}

func (m *MiotTracingServImpl) handleData(message mttypes.Message) (err error) {
	var taxiData mttypes.TaxiData
	if err = json.Unmarshal(message.Content, &taxiData); err != nil {
		iotlog.Errorf("json unmarshal error: %v", err)
		return
	}
	//insert data
	err = m.dp.InsertTaxiData(&mttypes.TaxiData{
		TaxiID:    taxiData.TaxiID,
		Timestamp: taxiData.Timestamp,
		Longitude: taxiData.Longitude,
		Latitude:  taxiData.Latitude,
		Occupancy: taxiData.Occupancy,
		TaxiDataLabel: mttypes.TaxiDataLabel{
			Segment: taxiData.Segment,
		},
	})
	if err != nil {
		iotlog.Errorf("InsertTaxiData error: %v", err)
		return
	}
	return
}

// GetProcessChan 提供对 processChan 的访问
func (m *MiotTracingServImpl) GetProcessChan() chan mttypes.Message {
	return m.processChan
}

// GetSendingChan 提供对 sendingChan 的访问
func (m *MiotTracingServImpl) GetSendingChan() chan mttypes.Message {
	return m.sendingChan
}

// -------------------------helper function---------------------------
func decreaseSegment(segment string) string {
	segment_int, err := strconv.Atoi(segment)
	if err != nil {
		iotlog.Errorf("strconv.Atoi failed, err: %v", err)
		return ""
	}
	segment_int--
	return strconv.Itoa(segment_int)
}
