package secondindexprocessor

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

type SecondIndexProcessor struct {
	c *redis.Client
}

func NewSecondIndexProcessor(c *redis.Client) *SecondIndexProcessor {
	return &SecondIndexProcessor{c}
}

// index生成样子（3byte）：binNum(5位) + elementCode(11位) + segment(8位)
func (i *SecondIndexProcessor) CreateSecondIndex(mtdt *mttypes.SecondIndex) error {
	// 先把数据序列化
	XYT_compressed, err := compressXYT(mtdt.StartTs)
	if err != nil {
		iotlog.Errorln("compressXYT failed, err:", err)
		return err
	}
	segment, err := strconv.ParseInt(mtdt.Segment, 10, 64)
	if err != nil {
		iotlog.Errorln("strconv.Atoi failed, err:", err)
		return err
	}
	add_segment_index := combineXYTAndSegment(XYT_compressed, segment)
	RedisListKey := fmt.Sprintf("%s%s:%s:%s", mttypes.Node_prefix, mttypes.NODE_ID, mttypes.SecondIndex_prefix, mtdt.ID)
	ic := i.c.RPush(RedisListKey, add_segment_index)
	if ic.Err() != nil {
		iotlog.Errorln("RPush failed, err:", ic.Err())
		return ic.Err()
	}
	return nil
}

// 有个问题是 如果移动终端到达一个新节点后立刻又掉头回到原本节点 那么节点的元数据的最后一个
// 补充成完整的secondindex
func (i *SecondIndexProcessor) UpdateSecondIndex(mtdt *mttypes.SecondIndex) error {
	//从右边开始寻找
	RedisListKey := fmt.Sprintf("%s%s:%s:%s", mttypes.Node_prefix, mttypes.NODE_ID, mttypes.SecondIndex_prefix, mtdt.ID)
	sc := i.c.RPop(RedisListKey)
	if sc.Err() != nil {
		iotlog.Errorln("RPop failed, err:", sc.Err())
		return sc.Err()
	}

	buf, err := sc.Int64()
	if err != nil {
		iotlog.Errorln("Bytes failed, err:", err)
		return err
	}

	//拆分出segment和XYT
	segment, XYT := splitSegmentAndXYT(buf)

	//如果segment不相等，报错
	segment_form_message, err := strconv.ParseInt(mtdt.Segment, 10, 64)
	if err != nil {
		iotlog.Errorln("strconv.Atoi failed, err:", err)
		return err
	}
	if segment != segment_form_message {
		iotlog.Errorln("segment not equal")
		return fmt.Errorf("segment not equal")
	}

	//增加endtime和nextnode
	XYT_endtime_compressed, err := compressXYT(mtdt.EndTs)
	if err != nil {
		iotlog.Errorln("compressXYT failed, err:", err)
		return err
	}

	//nextnode
	nextnode, err := strconv.ParseInt(mtdt.NextNode, 10, 64)
	if err != nil {
		iotlog.Errorln("strconv.Atoi failed, err:", err)
		return err
	}
	//合并
	combined_whole_index := combineAll(combineStartXYTAndEndXYT(XYT, XYT_endtime_compressed), nextnode, segment)
	//把合并后的元素push回去
	i.c.RPush(RedisListKey, combined_whole_index)
	return nil
}

// 获取符合条件的secondindex
func (i *SecondIndexProcessor) getSecondIndex(id string, startTs_from_query string, endTs_from_query string) (second_indexes []mttypes.SecondIndex, err error) {
	RedisListKey := fmt.Sprintf("%s%s:%s:%s", mttypes.Node_prefix, mttypes.NODE_ID, mttypes.SecondIndex_prefix, id)
	//从左到右查询
	indexes, err := i.c.LRange(RedisListKey, 0, -1).Result()
	if err != nil {
		iotlog.Errorln("Error retrieving list elements:", err)
		return
	}
	//遍历每个元素
	for _, v := range indexes {
		//把v转成int64
		value_64, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			iotlog.Errorln("strconv.Atoi failed, err:", err)
			return nil, err
		}

		//解压
		start_ts, end_ts, segment, next_node := decompressSecondIndex(value_64)
		if check := checkSecondIndex(start_ts, end_ts, segment, next_node); !check {
			iotlog.Errorln("secondindex checked failed")
			return nil, fmt.Errorf("secondindex checked failed")
		}
		//如果start_ts大于等于startTs，就把这个secondindex加入到index中
		start_ts_64, err := strconv.ParseInt(start_ts, 10, 64)
		if err != nil {
			iotlog.Errorln("strconv.Atoi failed, err:", err)
			return nil, err
		}

		startTs_from_query_64, err := strconv.ParseInt(startTs_from_query, 10, 64)
		if err != nil {
			iotlog.Errorln("strconv.Atoi failed, err:", err)
			return nil, err
		}
		endTs_from_query_64, err := strconv.ParseInt(endTs_from_query, 10, 64)
		if err != nil {
			iotlog.Errorln("strconv.Atoi failed, err:", err)
			return nil, err
		}
		if start_ts_64 >= startTs_from_query_64 && start_ts_64 <= endTs_from_query_64 {
			second_indexes = append(second_indexes, mttypes.SecondIndex{
				ID:       id,
				StartTs:  start_ts,
				EndTs:    end_ts,
				Segment:  segment,
				NextNode: next_node,
			})
		} else {
			return nil, nil
		}

	}
	return
}

//-----------------helper functions-----------------

func checkSecondIndex(start string, end string, segment string, next_node string) bool {
	if start == "" || segment == "" || next_node == "" {
		return false
	} else if start == end {
		return false
	}
	return true
}
