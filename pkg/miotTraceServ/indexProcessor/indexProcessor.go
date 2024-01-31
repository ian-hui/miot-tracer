package indexprocessor

import (
	"encoding/json"
	"fmt"
	mttypes "miot_tracing_go/mtTypes"
	"miot_tracing_go/pkg/logger"
	secondindexprocessor "miot_tracing_go/pkg/miotTraceServ/indexProcessor/SecondIndexProcessor"
	"strconv"

	"github.com/go-redis/redis"
)

var (
	iotlog = logger.Miotlogger
	redisC *redis.Client
)

type IndexProcessor struct {
	c   *redis.Client
	SIP *secondindexprocessor.SecondIndexProcessor
}

func NewIndexProcessor() *IndexProcessor {
	if redisC == nil {
		c := redis.NewClient(&redis.Options{
			Addr:     mttypes.RedisConfig.Addr,
			Password: mttypes.RedisConfig.Pwd,
			DB:       0,
		})
		redisC = c
	}
	return &IndexProcessor{redisC,
		secondindexprocessor.NewSecondIndexProcessor(redisC)}
}

func (i *IndexProcessor) CreateThirdIndex(index *mttypes.ThirdIndex) error {
	//序列化
	value_json, err := json.Marshal(index)
	if err != nil {
		iotlog.Errorln("json.Marshal failed, err:", err)
		return err
	}
	//把序列化后的数据存入redis
	redisKeyName := fmt.Sprintf("%s%s:%s:%s", mttypes.Node_prefix, mttypes.NODE_ID, mttypes.ThirdIndex_prefix, index.ID)
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

//-----------------helper functions-----------------
