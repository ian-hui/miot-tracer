package indexprocessor

import (
	mttypes "miot_tracing_go/mtTypes"
	"miot_tracing_go/pkg/logger"
	secondindexprocessor "miot_tracing_go/pkg/miotTraceServ/indexProcessor/SecondIndexProcessor"
	thirdindexprocessor "miot_tracing_go/pkg/miotTraceServ/indexProcessor/ThirdIndexProcessor"

	"github.com/go-redis/redis"
)

var (
	iotlog = logger.Miotlogger
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

//-----------------helper functions-----------------
