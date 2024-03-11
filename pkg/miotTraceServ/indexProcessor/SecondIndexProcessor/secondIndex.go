package secondindexprocessor

import (
	"fmt"
	mttypes "miot_tracing_go/mtTypes"
	"miot_tracing_go/pkg/logger"
	"strconv"
	"strings"
	"time"

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

// index生成样子（4byte）：XYT(16bit) + segment(16位)
func (i *SecondIndexProcessor) CreateSecondIndex(mtdt *mttypes.SecondIndex) (err error) {
	RedisListKey := fmt.Sprintf("%s%s:%s:%s", mttypes.Node_prefix, mttypes.NODE_ID, mttypes.SecondIndex_prefix, mtdt.ID)

	add_segment_index := mtdt.StartTs + ":" + mtdt.Segment
	// // 先把数据序列化
	// XYT_compressed, err := compressXYT(mtdt.StartTs, 11)
	// if err != nil {
	// 	iotlog.Errorln("compressXYT failed, err:", err)
	// 	return err
	// }

	// segment, err := strconv.ParseInt(mtdt.Segment, 10, 64)
	// if err != nil {
	// 	iotlog.Errorln("strconv.Atoi failed, err:", err)
	// 	return err
	// }

	// add_segment_index := combineXYTAndSegment(XYT_compressed, segment)
	redis_list := i.c.RPush(RedisListKey, add_segment_index)
	if redis_list.Err() != nil {
		iotlog.Errorln("RPush failed, err:", redis_list.Err())
		return redis_list.Err()
	}
	iotlog.Infoln("CreateSecondIndex success, segment:", mtdt.Segment, "start_ts:", mtdt.StartTs)
	return nil
}

// todo 增加segment在list中的index
func (i *SecondIndexProcessor) UpdateSecondIndex(mtdt *mttypes.SecondIndex) error {
	RedisListKey := fmt.Sprintf("%s%s:%s:%s", mttypes.Node_prefix, mttypes.NODE_ID, mttypes.SecondIndex_prefix, mtdt.ID)

	//如果一条都没有，报错
	RL_len := i.c.LLen(RedisListKey)
	if RL_len.Err() != nil {
		iotlog.Errorln("LLen failed, err:", RL_len.Err())
		return RL_len.Err()
	}
	if RL_len.Val() == 0 {
		iotlog.Errorln("RL_len is 0")
		return fmt.Errorf("RL_len is 0")
	}

	//从左到右查询
	indexes, err := i.c.LRange(RedisListKey, 0, -1).Result()
	if err != nil {
		iotlog.Errorln("Error retrieving list elements:", err)
		return err
	}
	for retry := 0; retry < mttypes.RETRY; retry++ {

		for index, v := range indexes {
			// int64_v, err := strconv.ParseInt(v, 10, 64)
			// if err != nil {
			// 	iotlog.Errorln("strconv.Atoi failed, err:", err)
			// 	return err
			// }
			// //拆分出segment和XYT
			// segment, XYT := splitSegmentAndXYT(int64_v)
			//如果没有:,continue
			if !strings.Contains(v, ":") {
				continue
			}
			startTs_and_segment := strings.Split(v, ":")
			if len(startTs_and_segment) != 2 {
				iotlog.Errorln("len(startTs_and_segment) != 2")
				return fmt.Errorf("len(startTs_and_segment) != 2")
			}
			startTs, segment := startTs_and_segment[0], startTs_and_segment[1]

			segment_int64, err := strconv.ParseInt(segment, 10, 64)
			if err != nil {
				iotlog.Errorln("strconv.Atoi failed, err:", err)
				return err
			}

			//如果segment不相等，报错
			segment_form_message, err := strconv.ParseInt(mtdt.Segment, 10, 64)
			if err != nil {
				iotlog.Errorln("strconv.Atoi failed, err:", err)
				return err
			}
			//找到segment_form_message对应的index
			if segment_int64 != segment_form_message {
				continue
			}

			XYT_start_compressed, err := compressXYT(mtdt.StartTs, 11)
			if err != nil {
				iotlog.Errorln("compressXYT failed, err:", err)
				return err
			}

			//variableIndex
			XYT_endtime_compressed, err := VariableLengthCompress(mtdt.EndTs, startTs)
			if err != nil {
				iotlog.Errorln("VariableLengthCompress failed, err:", err)
				return err
			}

			//nextnode
			nextnode, err := strconv.ParseInt(mtdt.NextNode, 10, 64)
			if err != nil {
				iotlog.Errorln("strconv.Atoi failed, err:", err)
				return err
			}
			//合并
			combined_whole_index := combineAll(combineStartXYTAndEndXYT(XYT_start_compressed, XYT_endtime_compressed), segment_int64, nextnode)
			//把合并后的元素放回对应的位置
			sc := i.c.LSet(RedisListKey, int64(index), combined_whole_index)
			if sc.Err() != nil {
				iotlog.Errorln("LSet failed, err:", sc.Err())
				return sc.Err()
			}
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("segment not found segment_form_message:%s", mtdt.Segment)
}

// 获取符合条件的secondindex
func (i *SecondIndexProcessor) GetSecondIndex(id string, startTs_from_query string, endTs_from_query string) (second_indexes []mttypes.SecondIndex, err error) {
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
		var start_ts, end_ts, segment, next_node = "", "", "", ""
		//解压
		start_ts, end_ts, segment, next_node, err = decompressSecondIndex(value_64)
		if err != nil {
			iotlog.Errorln("decompressSecondIndex failed, err:", err)
			return nil, err
		}
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
		end_ts_64, err := strconv.ParseInt(end_ts, 10, 64)
		if err != nil {
			iotlog.Errorln("strconv.Atoi failed, err:", err)
			return nil, err
		}
		// fmt.Println(time.Unix(start_ts_64, 0).UTC())
		// end_ts_64, _ := strconv.ParseInt(end_ts, 10, 64)
		// fmt.Println(time.Unix(end_ts_64, 0).UTC())

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
		iotlog.Infoln("segment", segment, "start_ts_64:", start_ts_64, "end_ts_64:", end_ts_64, "startTs_from_query_64:", startTs_from_query_64, "endTs_from_query_64:", endTs_from_query_64)
		//只要不是startts和endts都在范围外的都加入
		if start_ts_64 <= endTs_from_query_64 && end_ts_64 >= startTs_from_query_64 {
			second_indexes = append(second_indexes, mttypes.SecondIndex{
				ID:       id,
				StartTs:  start_ts,
				EndTs:    end_ts,
				Segment:  segment,
				NextNode: next_node,
			})
		}
	}
	return
}

func (i *SecondIndexProcessor) FindNearestSegment(id string, startTs_from_query string) (find bool, close_or_finded_time mttypes.SecondIndex, err error) {
	RedisListKey := fmt.Sprintf("%s%s:%s:%s", mttypes.Node_prefix, mttypes.NODE_ID, mttypes.SecondIndex_prefix, id)
	close_or_finded_time, find = mttypes.SecondIndex{}, false

	//从左到右查询
	indexes, err := i.c.LRange(RedisListKey, 0, -1).Result()
	if err != nil {
		iotlog.Errorln("Error retrieving list elements:", err)
		return
	}
	startTs_from_query_64, err := strconv.ParseInt(startTs_from_query, 10, 64)
	if err != nil {
		iotlog.Errorln("strconv.Atoi failed, err:", err)
		return
	}
	//遍历每个元素
	for _, v := range indexes {
		//把v转成int64
		value_64, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			iotlog.Errorln("strconv.Atoi failed, err:", err)
			return false, mttypes.SecondIndex{}, err
		}
		var start_ts, end_ts, segment, next_node = "", "", "", ""
		//解压
		start_ts, end_ts, segment, next_node, err = decompressSecondIndex(value_64)
		if err != nil {
			iotlog.Errorln("decompressSecondIndex failed, err:", err)
			return false, mttypes.SecondIndex{}, err
		}
		if check := checkSecondIndex(start_ts, end_ts, segment, next_node); !check {
			iotlog.Errorln("secondindex checked failed")
			return false, mttypes.SecondIndex{}, fmt.Errorf("secondindex checked failed")
		}
		start_ts_64, err := strconv.ParseInt(start_ts, 10, 64)
		if err != nil {
			iotlog.Errorln("strconv.Atoi failed, err:", err)
			return false, mttypes.SecondIndex{}, err
		}
		end_ts_64, err := strconv.ParseInt(end_ts, 10, 64)
		if err != nil {
			iotlog.Errorln("strconv.Atoi failed, err:", err)
			return false, mttypes.SecondIndex{}, err
		}
		if start_ts_64 <= startTs_from_query_64 && end_ts_64 >= startTs_from_query_64 {
			find = true
			close_or_finded_time = mttypes.SecondIndex{
				ID:       id,
				StartTs:  start_ts,
				EndTs:    end_ts,
				Segment:  segment,
				NextNode: next_node,
			}
		} else if start_ts_64 < startTs_from_query_64 {
			close_or_finded_time = mttypes.SecondIndex{
				ID:       id,
				StartTs:  start_ts,
				EndTs:    end_ts,
				Segment:  segment,
				NextNode: next_node,
			}
		} else {
			continue
		}
	}
	return find, close_or_finded_time, nil
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
