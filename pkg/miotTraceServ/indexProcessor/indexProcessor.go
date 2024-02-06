package indexprocessor

import (
	mttypes "miot_tracing_go/mtTypes"
	secondindexprocessor "miot_tracing_go/pkg/miotTraceServ/indexProcessor/SecondIndexProcessor"
	thirdindexprocessor "miot_tracing_go/pkg/miotTraceServ/indexProcessor/ThirdIndexProcessor"
	"strconv"

	"github.com/go-redis/redis"
)

var (
	redisC *redis.Client
)

type IndexProcessor struct {
	c   *redis.Client
	SIP *secondindexprocessor.SecondIndexProcessor
	TIP *thirdindexprocessor.ThirdIndexProcessor
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
		secondindexprocessor.NewSecondIndexProcessor(redisC),
		thirdindexprocessor.NewThirdIndexProcessor(redisC),
	}
}

func (ip *IndexProcessor) InsertHeadMeta(f_data mttypes.FirstData) error {
	meta := f_data.Timestamp + ":" + strconv.Itoa(mttypes.TS_SKIP)
	sc := ip.c.Set(f_data.TaxiID, meta, 0)
	return sc.Err()
}
