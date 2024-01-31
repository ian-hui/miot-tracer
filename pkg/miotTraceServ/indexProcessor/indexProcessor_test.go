package indexprocessor

import (
	"fmt"
	mttypes "miot_tracing_go/mtTypes"
	"strconv"
	"testing"
	"time"
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
	err := i.CreateSecondIndex(&test0)
	if err != nil {
		t.Error(err)
	}
	err2 := i.UpdateSecondIndex(&test)
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

func TestXYT(t *testing.T) {
	layout := "2006-01-02 15:04:05"
	times, err := time.Parse(layout, "2008-01-02 12:30:57")
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}
	s := strconv.FormatInt(times.Unix(), 10)
	fmt.Println(s)
	combined := compressXYT(s)
	fmt.Println(combined)
	fmt.Println(decompressXYT(combined))
}
