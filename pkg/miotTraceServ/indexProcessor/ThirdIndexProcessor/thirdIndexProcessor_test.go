package thirdindexprocessor

import (
	"fmt"
	mttypes "miot_tracing_go/mtTypes"
	"testing"

	"github.com/go-redis/redis"
)

var (
	s = NewThirdIndexProcessor(redis.NewClient(&redis.Options{
		Addr:     mttypes.RedisConfig.Addr,
		Password: mttypes.RedisConfig.Pwd,
		DB:       0,
	}))
)

func TestAddThirdIndex(t *testing.T) {
	var (
		testIndex = mttypes.ThirdIndex{
			ID:          "1",
			SequenceNum: "2",
			NodeID:      "3",
		}
		testQuery = mttypes.QueryStru{
			ID: "1",
			Tii: mttypes.ThirdIndexInfo{
				Skip_Ts:       "1800",
				Taxi_Start_Ts: "1800",
			},
			StartTime: "10000",
		}
	)
	err := s.CreateThirdIndex(&testIndex)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}
	node_id, err := s.QueryThirdIndex(&testQuery)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}
	fmt.Println(node_id)
}
