package thirdindexprocessor

import (
	"fmt"
	mttypes "miot_tracing_go/mtTypes"
	"miot_tracing_go/pkg/logger"
	"strconv"

	"github.com/go-redis/redis"
)

var (
	iotlog = logger.Miotlogger
)

type ThirdIndexProcessor struct {
	c *redis.Client
}

func NewThirdIndexProcessor(c *redis.Client) *ThirdIndexProcessor {
	return &ThirdIndexProcessor{c}
}

func (t *ThirdIndexProcessor) CreateThirdIndex(index *mttypes.ThirdIndex) error {
	redisKeyName := fmt.Sprintf("%s%s:%s:%s", mttypes.Node_prefix, mttypes.NODE_ID, mttypes.ThirdIndex_prefix, index.ID)
	fmt.Println(redisKeyName)
	compressed, err := compressThirdIndex(index.SequenceNum, index.NodeID)
	if err != nil {
		iotlog.Errorln("compressThirdIndex failed, err:", err)
		return err
	}
	float64_compressed := float64(compressed)
	string_compressed := strconv.FormatInt(compressed, 10)
	//利用timestamp作为score，这样子搜索的时候就可以用zrangebyscore
	t.c.ZAdd(redisKeyName, redis.Z{Score: float64_compressed, Member: string_compressed})
	return nil
}

func (t *ThirdIndexProcessor) QueryThirdIndex(query *mttypes.QueryStru) (node_id string, err error) {
	redisKeyName := fmt.Sprintf("%s%s:%s:%s", mttypes.Node_prefix, mttypes.NODE_ID, mttypes.ThirdIndex_prefix, query.ID)
	//从分数高开始遍历
	ssc := t.c.ZRevRange(redisKeyName, 0, -1)
	if ssc.Err() != nil {
		iotlog.Errorln("ZRevRange failed, err:", ssc.Err())
		return "", ssc.Err()
	}
	for _, v := range ssc.Val() {
		//解压
		combined, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			iotlog.Errorln("strconv.ParseInt failed, err:", err)
			return "", err
		}
		sequence_num, node_id, err := decompressThirdIndex(combined)
		if err != nil {
			iotlog.Errorln("decompressThirdIndex failed, err:", err)
			return "", err
		}
		query_skip_ts_64, err := strconv.ParseInt(query.Tii.Skip_Ts, 10, 64)
		if err != nil {
			iotlog.Errorln("strconv.ParseInt failed, err:", err)
			return "", err
		}
		query_taxi_start_ts_64, err := strconv.ParseInt(query.Tii.Taxi_Start_Ts, 10, 64)
		if err != nil {
			iotlog.Errorln("strconv.ParseInt failed, err:", err)
			return "", err
		}
		sequence_num_64, err := strconv.ParseInt(sequence_num, 10, 64)
		if err != nil {
			iotlog.Errorln("strconv.ParseInt failed, err:", err)
			return "", err
		}
		query_startTime_64, err := strconv.ParseInt(query.StartTime, 10, 64)
		if err != nil {
			iotlog.Errorln("strconv.ParseInt failed, err:", err)
			return "", err
		}
		//找到第一个满足条件的
		if sequence_num_64*query_skip_ts_64+query_taxi_start_ts_64 <= query_startTime_64 {
			return node_id, nil
		}
	}
	return "", nil
}

// 这个float64是一个组合数，前面是sequence_num，后面是node_id
func compressThirdIndex(sequence_num string, node_id string) (int64, error) {
	//序列号
	sequence_num_int64, err := strconv.ParseInt(sequence_num, 10, 64)
	if err != nil {
		iotlog.Errorln("strconv.ParseInt failed, err:", err)
		return 0, err
	}
	//节点id
	node_id_int64, err := strconv.ParseInt(node_id, 10, 64)
	if err != nil {
		iotlog.Errorln("strconv.ParseInt failed, err:", err)
		return 0, err
	}
	//组合
	combined := (sequence_num_int64 << mttypes.THIRD_INDEX_NODE_ID_LEN) | node_id_int64
	return combined, nil
}

func decompressThirdIndex(combined int64) (string, string, error) {
	sequence_num := combined >> mttypes.THIRD_INDEX_NODE_ID_LEN
	node_id := combined & (1<<mttypes.THIRD_INDEX_NODE_ID_LEN - 1)
	return strconv.FormatInt(sequence_num, 10), strconv.FormatInt(node_id, 10), nil
}
