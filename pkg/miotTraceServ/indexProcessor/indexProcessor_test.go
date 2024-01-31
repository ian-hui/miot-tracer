package indexprocessor

import (
	mttypes "miot_tracing_go/mtTypes"
	"testing"
)

var (
	test0 = mttypes.SecondIndex{
		ID:      "1",
		StartTs: "1",
		Segment: "3",
	}
	test = mttypes.SecondIndex{
		ID:       "1",
		EndTs:    "2",
		Segment:  "3",
		NextNode: "4",
	}
	testIndex = mttypes.ThirdIndex{
		ID:        "1",
		Timestamp: "2",
		NodeID:    "3",
		Segment:   "4",
	}
)

func TestAddMeta(t *testing.T) {

	i := NewIndexProcessor()
	err := i.SIP.CreateSecondIndex(&test0)
	if err != nil {
		t.Error(err)
	}
	err2 := i.SIP.UpdateSecondIndex(&test)
	if err2 != nil {
		t.Error(err2)
	}
}

func TestCreateIndex(t *testing.T) {

	i := NewIndexProcessor()
	err := i.CreateThirdIndex(&testIndex)
	if err != nil {
		t.Error(err)
	}
}
