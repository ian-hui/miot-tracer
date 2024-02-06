package miottraceserv

import (
	"encoding/json"
	"fmt"
	"math"
	mttypes "miot_tracing_go/mtTypes"
	consistenthash "miot_tracing_go/pkg/consistentHash"
	"miot_tracing_go/pkg/logger"
	dataprocessor "miot_tracing_go/pkg/miotTraceServ/dataProcessor"
	indexprocessor "miot_tracing_go/pkg/miotTraceServ/indexProcessor"
	"strconv"
	"strings"
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
	consistentHash *consistenthash.ConsistentHash
	dp             *dataprocessor.DataProcessor
	ip             *indexprocessor.IndexProcessor
	processChan    chan mttypes.Message
	sendingChan    chan mttypes.Message
}

func NewMiotTracingServImpl() MiotTracingServ {
	return &MiotTracingServImpl{
		consistentHash: consistenthash.NewConsistentHash(10),
		dp:             dataprocessor.NewDataProcessor(),
		ip:             indexprocessor.NewIndexProcessor(),
		processChan:    make(chan mttypes.Message, 100),
		sendingChan:    make(chan mttypes.Message, 100),
	}
}

func (mts *MiotTracingServImpl) Start(worker_num int) {
	// 初始化一致性哈希
	for i := 1; i <= mttypes.NODE_TOTAL_NUM; i++ {
		mts.consistentHash.AddNode(fmt.Sprintf("%d", i))
	}
	// 启动worker
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
		err := m.handleFirstData(message)
		if err != nil {
			iotlog.Errorf("handleFirstData error: %v", err)
		}
	case mttypes.TYPE_UPLOAD_THIRD_INDEX:
		err := m.handleUploadThirdIndex(message)
		if err != nil {
			iotlog.Errorf("handleUploadThirdIndex error: %v", err)
		}
	case mttypes.TYPE_UPLOAD_TAXI_DATA:
		err := m.handleData(message)
		if err != nil {
			iotlog.Errorf("handleData error: %v", err)
		}
	case mttypes.TYPE_UPDATE_SECOND_INDEX:
		err := m.handleUpdateSecondIndex(message)
		if err != nil {
			iotlog.Errorf("handleUpdateSecondIndex error: %v", err)
		}
	case mttypes.TYPE_UPDATE_THIRD_INDEX:
		err := m.handleUpdateThirdIndex(message)
		if err != nil {
			iotlog.Errorf("handleUpdateThirdIndex error: %v", err)
		}
	case mttypes.TYPE_UPLOAD_THIRD_INDEX_HEAD:
		err := m.handleUploadMetaData(message)
		if err != nil {
			iotlog.Errorf("handleUploadMetaData error: %v", err)
		}
	default:
		iotlog.Errorf("unknown message type: %v", message.Type)
	}
	return nil
}

func (m *MiotTracingServImpl) handleFirstData(message mttypes.Message) (err error) {
	var firstData mttypes.FirstData
	if err = json.Unmarshal(message.Content, &firstData); err != nil {
		iotlog.Errorf("json unmarshal error: %v", err)
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

		if firstData.TaxiFrontNode == "" {
			iotlog.Errorf("empty TaxiFrontNode")
			return
		}
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
		topic := string(firstData.TaxiFrontNode)
		m.sendingChan <- mttypes.Message{
			Type:    mttypes.TYPE_UPDATE_SECOND_INDEX,
			Topic:   topic,
			Content: sendback_second_index_json,
		}
	} else {
		//第一次上传数据，初始化
		//传递信息-上传链表头节点（third index header）
		node_id := m.consistentHash.GetNode(mttypes.TAXI_ID_PREFIX + ":" + firstData.TaxiID)
		if node_id != mttypes.NODE_ID {
			topic := node_id
			m.sendingChan <- mttypes.Message{
				Type:    mttypes.TYPE_UPLOAD_THIRD_INDEX_HEAD,
				Topic:   topic,
				Content: message.Content,
			}
		} else {
			err = m.ip.InsertHeadMeta(firstData)
			if err != nil {
				iotlog.Errorf("InsertHeadMeta error: %v", err)
				return
			}
		}
		//更新第三级索引
		third_index := mttypes.ThirdIndex{
			ID:          firstData.TaxiID,
			SequenceNum: "1",
			NodeID:      mttypes.NODE_ID,
		}
		third_index_json, err := json.Marshal(third_index)
		if err != nil {
			iotlog.Errorf("json marshal error: %v", err)
			return err
		}
		topic := node_id
		m.sendingChan <- mttypes.Message{
			Type:    mttypes.TYPE_UPDATE_THIRD_INDEX,
			Topic:   topic,
			Content: third_index_json,
		}
	}
	return
}

func (m *MiotTracingServImpl) handleUploadMetaData(message mttypes.Message) (err error) {
	var firstData mttypes.FirstData
	if err = json.Unmarshal(message.Content, &firstData); err != nil {
		iotlog.Errorf("json unmarshal error: %v", err)
		return
	}
	err = m.ip.InsertHeadMeta(firstData)
	if err != nil {
		iotlog.Errorf("InsertHeadMeta error: %v", err)
		return
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

// 收到这个说明现在这个节点就是index的最后一个节点了
func (m *MiotTracingServImpl) handleUploadThirdIndex(message mttypes.Message) (err error) {
	var taxi_info mttypes.TaxiInfo
	if err = json.Unmarshal(message.Content, &taxi_info); err != nil {
		iotlog.Errorf("json unmarshal error: %v", err)
		return
	}
	//taxi记录的索引是没有头节点的，所以要加上头节点
	the_header_index := m.consistentHash.GetNode(mttypes.TAXI_ID_PREFIX + ":" + taxi_info.TaxiID)
	taxi_info.Index = append([]string{"0" + ":" + the_header_index}, taxi_info.Index...)
	forward_map, err := getForwardThirdIndexMap(taxi_info.Index)
	if err != nil {
		iotlog.Errorf("getForwardThirdIndexMap error: %v", err)
		return
	}
	for node_id, third_index := range forward_map {
		third_index_json, err := json.Marshal(third_index)
		if err != nil {
			iotlog.Errorf("json marshal error: %v", err)
			return err
		}
		topic := node_id
		m.sendingChan <- mttypes.Message{
			Type:    mttypes.TYPE_UPDATE_THIRD_INDEX,
			Topic:   topic,
			Content: third_index_json,
		}
	}
	return
}

func (m *MiotTracingServImpl) handleUpdateSecondIndex(message mttypes.Message) (err error) {
	var second_index mttypes.SecondIndex
	if err = json.Unmarshal(message.Content, &second_index); err != nil {
		iotlog.Errorf("json unmarshal error: %v", err)
		return
	}
	err = m.ip.SIP.UpdateSecondIndex(&second_index)
	if err != nil {
		iotlog.Errorf("UpdateSecondIndex error: %v", err)
		return
	}
	return
}

func (m *MiotTracingServImpl) handleUpdateThirdIndex(message mttypes.Message) (err error) {
	var third_index mttypes.ThirdIndex
	if err = json.Unmarshal(message.Content, &third_index); err != nil {
		iotlog.Errorf("json unmarshal error: %v", err)
		return
	}
	err = m.ip.TIP.CreateThirdIndex(&third_index)
	if err != nil {
		iotlog.Errorf("UpdateThirdIndex error: %v", err)
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

// seq:id
func getForwardThirdIndexMap(indexes []string) (forward_map map[string]mttypes.ThirdIndex, err error) {
	//取出本地的thirdindex
	local_index := len(indexes) - 1
	seq_and_id := strings.Split(indexes[local_index], ":")
	if len(seq_and_id) != 2 {
		iotlog.Errorf("invalid index: %v", indexes[local_index])
		return nil, fmt.Errorf("invalid index: %v", indexes[local_index])
	}
	sequence_num, node_id := seq_and_id[0], seq_and_id[1]
	//获取需要转发的索引
	sequence_num_int, err := strconv.Atoi(sequence_num)
	if err != nil {
		iotlog.Errorf("strconv.Atoi failed, err: %v", err)
		return nil, err
	}
	slice_indexes, err := getNeedForwardIndexList(sequence_num_int)
	if err != nil {
		iotlog.Errorf("getNeedForwardIndexList failed, err: %v", err)
		return nil, err
	}
	forward_map = make(map[string]mttypes.ThirdIndex)
	for _, index := range slice_indexes {
		seq_forward_node_id := strings.Split(indexes[index], ":")
		if len(seq_forward_node_id) != 2 {
			iotlog.Errorf("invalid index: %v", indexes[index])
			return nil, fmt.Errorf("invalid index: %v", indexes[index])
		}
		forward_node_id := seq_forward_node_id[1]
		//把自己的thirdindex转发给需要的节点
		forward_map[forward_node_id] = mttypes.ThirdIndex{
			SequenceNum: sequence_num,
			NodeID:      node_id,
		}
	}
	return
}

// 获取数组中需要转发的索引
func getNeedForwardIndexList(seq int) (slice_indexes []int, err error) {
	if seq <= 0 {
		return nil, fmt.Errorf("invalid sequence number: %v", seq)
	}
	n := 0.0
	slice_indexes = []int{seq - 1}
	for seq%2 == 0 {
		n++
		seq >>= 1 // 使用位右移代替除以2
		slice_indexes = append(slice_indexes, seq-int(math.Pow(2, n)))
	}
	return
}
