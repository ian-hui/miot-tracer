package indexprocessor

import (
	mttypes "miot_tracing_go/mtTypes"
	"testing"

	"github.com/go-redis/redis"
)

var (
	test0 = mttypes.Metadata{
		ID:      "1",
		StartTs: "1",
		Segment: "3",
	}
	test = mttypes.Metadata{
		ID:       "1",
		EndTs:    "2",
		Segment:  "3",
		NextNode: "4",
	}
	testIndex = mttypes.Index{
		ID:        "1",
		Timestamp: "2",
		NodeID:    "3",
		Segment:   "4",
	}
)

func TestAddMeta(t *testing.T) {
	c := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "reins5401",
		DB:       0,
	})
	i := NewIndexProcessor(c)
	err := i.createMetaData(&test0)
	if err != nil {
		t.Error(err)
	}
	err2 := i.updateMetaData(&test)
	if err2 != nil {
		t.Error(err2)
	}
}

func TestCreateIndex(t *testing.T) {
	c := redis.NewClient(&redis.Options{
		Addr:     mttypes.RedisConfig.Addr,
		Password: mttypes.RedisConfig.Pwd,
		DB:       0,
	})
	i := NewIndexProcessor(c)
	err := i.createIndex(&testIndex)
	if err != nil {
		t.Error(err)
	}
}
