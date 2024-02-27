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
	GetResultChan() chan mttypes.Result
}

type MiotTracingServImpl struct {
	consistentHash *consistenthash.ConsistentHash
	dp             *dataprocessor.DataProcessor
	ip             *indexprocessor.IndexProcessor
	processChan    chan mttypes.Message
	sendingChan    chan mttypes.Message
	resultChan     chan mttypes.Result
}

func NewMiotTracingServImpl() MiotTracingServ {
	return &MiotTracingServImpl{
		consistentHash: consistenthash.NewConsistentHash(10),
		dp:             dataprocessor.NewDataProcessor(),
		ip:             indexprocessor.NewIndexProcessor(),
		processChan:    make(chan mttypes.Message, 100),
		sendingChan:    make(chan mttypes.Message, 100),
		resultChan:     make(chan mttypes.Result, 100),
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

func (m *MiotTracingServImpl) handleFirstData(message mttypes.Message) (err error) {
	contentBytes, err := json.Marshal(message.Content)
	if err != nil {
		iotlog.Errorln("Error marshalling content back to JSON:", err)
		return
	}
	var firstData mttypes.FirstData
	if err = json.Unmarshal(contentBytes, &firstData); err != nil {
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

	//是第一次在本节点上传数据，但是不是第一次上传数据
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
		topic := string(firstData.TaxiFrontNode)
		m.sendingChan <- mttypes.Message{
			Type:    mttypes.TYPE_UPDATE_SECOND_INDEX,
			Topic:   topic,
			Content: sendback_second_index,
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

		topic := node_id
		m.sendingChan <- mttypes.Message{
			Type:    mttypes.TYPE_UPDATE_THIRD_INDEX,
			Topic:   topic,
			Content: third_index,
		}
	}
	return
}

func (m *MiotTracingServImpl) handleUploadMetaData(message mttypes.Message) (err error) {
	contentBytes, err := json.Marshal(message.Content)
	if err != nil {
		iotlog.Errorln("Error marshalling content back to JSON:", err)
		return
	}
	var firstData mttypes.FirstData
	if err = json.Unmarshal(contentBytes, &firstData); err != nil {
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
	contentBytes, err := json.Marshal(message.Content)
	if err != nil {
		iotlog.Errorln("Error marshalling content back to JSON:", err)
		return
	}
	var taxiData mttypes.TaxiData
	if err = json.Unmarshal(contentBytes, &taxiData); err != nil {
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
	contentBytes, err := json.Marshal(message.Content)
	if err != nil {
		iotlog.Errorln("Error marshalling content back to JSON:", err)
		return
	}
	var taxi_info mttypes.TaxiInfo
	if err = json.Unmarshal(contentBytes, &taxi_info); err != nil {
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
	for node_id, third_indexs := range forward_map {
		topic := node_id
		for _, third_index := range third_indexs {

			third_index.ID = taxi_info.TaxiID //补充id

			m.sendingChan <- mttypes.Message{
				Type:    mttypes.TYPE_UPDATE_THIRD_INDEX,
				Topic:   topic,
				Content: third_index,
			}
		}
	}
	return
}

func (m *MiotTracingServImpl) handleUpdateSecondIndex(message mttypes.Message) (err error) {
	contentBytes, err := json.Marshal(message.Content)
	if err != nil {
		iotlog.Errorln("Error marshalling content back to JSON:", err)
		return
	}
	var second_index mttypes.SecondIndex
	if err = json.Unmarshal(contentBytes, &second_index); err != nil {
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
	contentBytes, err := json.Marshal(message.Content)
	if err != nil {
		iotlog.Errorln("Error marshalling content back to JSON:", err)
		return
	}
	var third_index mttypes.ThirdIndex
	if err = json.Unmarshal(contentBytes, &third_index); err != nil {
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

// 查询
func (m *MiotTracingServImpl) handleBuildQuery(message mttypes.Message) (err error) {
	contentBytes, err := json.Marshal(message.Content)
	if err != nil {
		iotlog.Errorln("Error marshalling content back to JSON:", err)
		return
	}
	var query mttypes.QueryStru
	if err = json.Unmarshal(contentBytes, &query); err != nil {
		iotlog.Errorf("json unmarshal error: %v", err)
		return
	}
	head_node_id := m.consistentHash.GetNode(mttypes.TAXI_ID_PREFIX + ":" + query.ID)
	if message.Topic == "" {
		//补全信息 计算出头节点
		topic := head_node_id
		query.QueryNode = mttypes.NODE_ID
		m.sendingChan <- mttypes.Message{
			Type:    mttypes.TYPE_BUILD_QUERY,
			Topic:   topic,
			Content: query,
		}
		return
	} else {
		//头节点操作
		//补全tti
		ts, ts_skip, err := m.ip.QueryHeadMeta(query.ID)
		if err != nil {
			iotlog.Errorf("QueryHeadMeta error: %v", err)
			return err
		}
		query.Tii.Taxi_Start_Ts = ts
		query.Tii.Skip_Ts = ts_skip
		//查询第三级索引
		node_id, err := m.ip.TIP.QueryThirdIndex(&query)
		if err != nil {
			iotlog.Errorf("QueryThirdIndex error: %v", err)
			return err
		}

		if node_id == "" || node_id == mttypes.NODE_ID {
			//开始执行查询
			m.processChan <- mttypes.Message{
				Type:    mttypes.TYPE_FIND_STARTER,
				Topic:   mttypes.NODE_ID,
				Content: query,
			}
		} else {
			//转发
			topic := node_id
			m.sendingChan <- mttypes.Message{
				Type:    mttypes.TYPE_SEARCH_THIRD_INDEX,
				Topic:   topic,
				Content: query,
			}
		}
	}
	return
}

func (m *MiotTracingServImpl) handleSearchThirdIndex(message mttypes.Message) (err error) {
	contentBytes, err := json.Marshal(message.Content)
	if err != nil {
		iotlog.Errorln("Error marshalling content back to JSON:", err)
		return
	}
	var query mttypes.QueryStru
	if err = json.Unmarshal(contentBytes, &query); err != nil {
		iotlog.Errorf("json unmarshal error: %v", err)
		return
	}
	node_id, err := m.ip.TIP.QueryThirdIndex(&query)
	if err != nil {
		iotlog.Errorf("QueryThirdIndex error: %v", err)
		return err
	}
	if node_id == "" || node_id == mttypes.NODE_ID {
		//开始执行查询
		iotlog.Infof("start query data", query)
		m.processChan <- mttypes.Message{
			Type:    mttypes.TYPE_FIND_STARTER,
			Topic:   mttypes.NODE_ID,
			Content: query,
		}
	} else {
		//转发
		topic := node_id
		m.sendingChan <- mttypes.Message{
			Type:    mttypes.TYPE_SEARCH_THIRD_INDEX,
			Topic:   topic,
			Content: query,
		}
	}
	return
}

func (m *MiotTracingServImpl) handleQueryData(message mttypes.Message) (err error) {
	contentBytes, err := json.Marshal(message.Content)
	if err != nil {
		iotlog.Errorln("Error marshalling content back to JSON:", err)
		return
	}
	var query mttypes.QueryStru
	if err = json.Unmarshal(contentBytes, &query); err != nil {
		iotlog.Errorf("json unmarshal error: %v", err)
		return
	}
	if_traversed := checkAlreadyTraversed(mttypes.NODE_ID, query.TraverseCfg.Traversed)
	if if_traversed {
		iotlog.Info("already traversed")
		return
	}
	second_indexes, err := m.ip.SIP.GetSecondIndex(query.ID, query.StartTime, query.EndTime)
	if second_indexes == nil {
		//没有数据

		m.sendingChan <- mttypes.Message{
			Type:  mttypes.TYPE_SEND_BACK_RESULT,
			Topic: query.QueryNode,
			Content: mttypes.Result{
				Query_node:  query.QueryNode,
				Request_id:  query.RequestID,
				End_segment: query.TraverseCfg.Previous_segment,
			},
		}

		return
	}
	//有2index
	for _, second_index := range second_indexes {
		var res mttypes.Result
		res.Request_id = query.RequestID
		res.Query_node = query.QueryNode
		data_res, err := m.dp.QueryTaxiData(&mttypes.QueryStru{
			ID:        query.ID,
			Segment:   second_index.Segment,
			StartTime: second_index.StartTs,
			EndTime:   second_index.EndTs,
		})
		if err != nil {
			iotlog.Errorf("QueryTaxiData error: %v", err)
			return err
		}
		if data_res == nil {
			iotlog.Errorf("QueryTaxiData error: %v", err)
			return err
		}
		//开始数据
		if query.TraverseCfg.Start_segment == "" || second_index.Segment < query.TraverseCfg.Start_segment {
			res.Start_segment = second_index.Segment
			query.TraverseCfg.Start_segment = second_index.Segment
		}
		//最后数据
		if second_index.NextNode == "" {
			res.End_segment = second_index.Segment
		}
		res.Result = data_res
		//query 配置
		query.TraverseCfg.Previous_segment = second_index.Segment
		query.TraverseCfg.Traversed = append(query.TraverseCfg.Traversed, mttypes.NODE_ID)
		//如果没有到最后一个节点，继续查询
		if res.End_segment == "" {
			m.sendingChan <- mttypes.Message{
				Type:    mttypes.TYPE_QUERY_TAXI_DATA,
				Topic:   second_index.NextNode,
				Content: query,
			}
		}
		//返回结果
		m.sendingChan <- mttypes.Message{
			Type:    mttypes.TYPE_SEND_BACK_RESULT,
			Topic:   query.QueryNode,
			Content: res,
		}
	}
	return
}

func (m *MiotTracingServImpl) handleFindStarterOfQuery(message mttypes.Message) (err error) {
	contentBytes, err := json.Marshal(message.Content)
	if err != nil {
		iotlog.Errorln("Error marshalling content back to JSON:", err)
		return
	}
	var query mttypes.QueryStru
	if err = json.Unmarshal(contentBytes, &query); err != nil {
		iotlog.Errorf("json unmarshal error: %v", err)
		return
	}
	find, si, err := m.ip.SIP.FindNearestSegment(query.ID, query.StartTime)
	if err != nil {
		iotlog.Errorf("FindNearestSegment error: %v", err)
		return
	}
	if find {
		m.processChan <- mttypes.Message{
			Type:    mttypes.TYPE_QUERY_TAXI_DATA,
			Topic:   mttypes.NODE_ID,
			Content: query,
		}
	} else {
		m.sendingChan <- mttypes.Message{
			Type:    mttypes.TYPE_FIND_STARTER,
			Topic:   si.NextNode,
			Content: query,
		}
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

// GetResultChan 提供对 resultChan 的访问
func (m *MiotTracingServImpl) GetResultChan() chan mttypes.Result {
	return m.resultChan
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

// seq:id //在获取3级索引的时候，需要把自己的3级索引转发给需要的节点
func getForwardThirdIndexMap(indexes []string) (forward_map map[string][]mttypes.ThirdIndex, err error) {
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
	forward_map = make(map[string][]mttypes.ThirdIndex)
	for _, index := range slice_indexes {
		seq_forward_node_id := strings.Split(indexes[index], ":")
		if len(seq_forward_node_id) != 2 {
			iotlog.Errorf("invalid index: %v", indexes[index])
			return nil, fmt.Errorf("invalid index: %v", indexes[index])
		}
		forward_node_id := seq_forward_node_id[1]
		//把自己的thirdindex转发给需要的节点
		forward_map[forward_node_id] = append(forward_map[forward_node_id], mttypes.ThirdIndex{
			SequenceNum: sequence_num,
			NodeID:      node_id,
		})
	}
	return
}

// 获取数组中需要转发的索引
func getNeedForwardIndexList(seq int) (slice_indexes []int, err error) {
	if seq <= 0 {
		return nil, fmt.Errorf("invalid sequence number: %v", seq)
	}
	n, origin_index := 0.0, seq
	slice_indexes = []int{seq - 1}
	for seq%2 == 0 {
		n++
		seq >>= 1 // 使用位右移代替除以2
		need_forward_node_id := origin_index - int(math.Pow(2, n))
		slice_indexes = append(slice_indexes, need_forward_node_id)
	}
	return
}

func checkAlreadyTraversed(node_id string, traversed []string) bool {
	for _, id := range traversed {
		if id == node_id {
			return true
		}
	}
	return false
}
