package indexprocessor

import (
	"encoding/json"
	"fmt"
	mttypes "miot_tracing_go/mtTypes"
	"miot_tracing_go/pkg/logger"
	"strconv"

	"github.com/go-redis/redis"
)

var (
	iotlog = logger.Miotlogger
)

type IndexProcessor interface {
	createMetaData(mtdt *mttypes.Metadata) error
	updateMetaData(mtdt *mttypes.Metadata) error
	createIndex(index *mttypes.Index) error
}

type IndexProcessorImpl struct {
	c *redis.Client
}

func NewIndexProcessor(c *redis.Client) IndexProcessor {
	return &IndexProcessorImpl{c}
}

func (i *IndexProcessorImpl) createMetaData(mtdt *mttypes.Metadata) error {
	value_json, err := json.Marshal(mtdt)
	if err != nil {
		iotlog.Errorln("json.Marshal failed, err:", err)
		return err
	}
	RedisListKey := fmt.Sprintf("%s%s:%s:%s", mttypes.Node_prefix, mttypes.NODE_ID, mttypes.Meta_prefix, mtdt.ID)
	ic := i.c.RPush(RedisListKey, value_json)
	if ic.Err() != nil {
		iotlog.Errorln("RPush failed, err:", ic.Err())
		return ic.Err()
	}
	return nil
}

// 有个问题是 如果移动终端到达一个新节点后立刻又掉头回到原本节点 那么节点的元数据的最后一个
func (i *IndexProcessorImpl) updateMetaData(mtdt *mttypes.Metadata) error {
	//从右边开始寻找
	RedisListKey := fmt.Sprintf("%s%s:%s:%s", mttypes.Node_prefix, mttypes.NODE_ID, mttypes.Meta_prefix, mtdt.ID)
	sc := i.c.RPop(RedisListKey)
	if sc.Err() != nil {
		iotlog.Errorln("RPop failed, err:", sc.Err())
		return sc.Err()
	}
	//把pop出来的元素unmarshal成struct
	b, err := sc.Bytes()
	if err != nil {
		iotlog.Errorln("Bytes failed, err:", err)
		return err
	}
	var metaBeforeCombination mttypes.Metadata
	err = json.Unmarshal(b, &metaBeforeCombination)
	if err != nil {
		iotlog.Errorln("json.Unmarshal failed, err:", err)
		return err
	}
	//把新的元素和pop出来的元素合并
	combineMetaData(&metaBeforeCombination, mtdt)
	//把合并后的元素push回去
	value_json, err := json.Marshal(mtdt)
	if err != nil {
		iotlog.Errorln("json.Marshal failed, err:", err)
		return err
	}
	i.c.RPush(RedisListKey, value_json)
	return nil
}

func combineMetaData(metaBeforeCombination *mttypes.Metadata, mtdt *mttypes.Metadata) *mttypes.Metadata {
	mtdt.StartTs = metaBeforeCombination.StartTs
	return metaBeforeCombination
}

func (i *IndexProcessorImpl) createIndex(index *mttypes.Index) error {
	//序列化
	value_json, err := json.Marshal(index)
	if err != nil {
		iotlog.Errorln("json.Marshal failed, err:", err)
		return err
	}
	//把序列化后的数据存入redis
	redisKeyName := fmt.Sprintf("%s%s:%s:%s", mttypes.Node_prefix, mttypes.NODE_ID, mttypes.Index_prefix, index.ID)
	fmt.Println(redisKeyName)
	float_ts, err := strconv.ParseFloat(index.Timestamp, 64)
	if err != nil {
		iotlog.Errorln("strconv.ParseFloat failed, err:", err)
		return err
	}
	//利用timestamp作为score，这样子搜索的时候就可以用zrangebyscore
	i.c.ZAdd(redisKeyName, redis.Z{Score: float_ts, Member: value_json})
	return nil
}
