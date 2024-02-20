package secondindexprocessor

import (
	"fmt"
	mttypes "miot_tracing_go/mtTypes"
	"strconv"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

var (
	layout = "2006-01-02 15:04:05"
	s      = NewSecondIndexProcessor(redis.NewClient(&redis.Options{
		Addr:     mttypes.RedisConfig.Addr,
		Password: mttypes.RedisConfig.Pwd,
		DB:       0,
	}))
)

func TestXYT(t *testing.T) {
	times, err := time.Parse(layout, "2008-01-06 09:01:01")
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}
	s := strconv.FormatInt(times.Unix(), 10)
	fmt.Println(s)
	//1199449857
	//1199610061
	combined, err := compressXYT(s, 10)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}
	// fmt.Println(combined)
	decompress_s := decompressXYT(combined, 10)
	fmt.Println(decompress_s)
	//转回时间戳
	i, _ := strconv.ParseInt(decompress_s, 10, 64)
	i2 := time.Unix(i, 0).UTC().Format(layout)
	fmt.Println(i2)
}

func TestCreate2Index(t *testing.T) {

	times, err := time.Parse(layout, "2008-01-02 12:30:57")
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}
	times_string := strconv.FormatInt(times.UTC().Unix(), 10)
	err = s.CreateSecondIndex(&mttypes.SecondIndex{
		ID:      "1",
		StartTs: times_string,
		Segment: "3",
	})
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}
	endtime, err := time.Parse(layout, "2008-01-05 16:15:58")
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}
	endtime_string := strconv.FormatInt(endtime.UTC().Unix(), 10)
	err = s.UpdateSecondIndex(&mttypes.SecondIndex{
		ID:       "1",
		EndTs:    endtime_string,
		Segment:  "3",
		NextNode: "4",
	})
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}

}

func TestAppend2Index(t *testing.T) {

	times, err := time.Parse(layout, "2008-01-06 09:01:01")
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}
	times_string := strconv.FormatInt(times.UTC().Unix(), 10)
	err = s.CreateSecondIndex(&mttypes.SecondIndex{
		ID:      "1",
		StartTs: times_string,
		Segment: "7",
	})
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}
	endtime, err := time.Parse(layout, "2008-01-07 21:11:03")
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}
	endtime_string := strconv.FormatInt(endtime.UTC().Unix(), 10)
	err = s.UpdateSecondIndex(&mttypes.SecondIndex{
		ID:       "1",
		EndTs:    endtime_string,
		Segment:  "7",
		NextNode: "4",
	})
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}

}

func TestQuery2Index(t *testing.T) {
	queryStart, err := time.Parse(layout, "2008-01-01 12:30:10")
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}
	queryStart_string := strconv.FormatInt(queryStart.UTC().Unix(), 10)
	queryEnd, err := time.Parse(layout, "2008-01-09 12:30:58")
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}
	queryEnd_string := strconv.FormatInt(queryEnd.UTC().Unix(), 10)

	second_indexes, err2 := s.GetSecondIndex("1", queryStart_string, queryEnd_string)
	if err2 != nil {
		fmt.Println("Error parsing date:", err2)
		return
	}
	fmt.Println(second_indexes)
}

func TestSplitAndCombineXYT(t *testing.T) {
	s := int64(3116)
	e := int64(5164)
	i := combineStartXYTAndEndXYT(s, e)
	fmt.Println(i)
	s1, e1 := splitXYT2StartEnd(i)
	fmt.Println(s1, e1)
	s2 := decompressXYT(s1, 11)
	e2 := decompressXYT(e1, 11)
	s2_64, _ := strconv.ParseInt(s2, 10, 64)
	e2_64, _ := strconv.ParseInt(e2, 10, 64)
	fmt.Println(time.Unix(s2_64, 0).UTC(), time.Unix(e2_64, 0).UTC())
	fmt.Println(s2, e2)
}

func TestSplitXYTAndSeg(t *testing.T) {
	segment, XYT := splitSegmentAndXYT(3178718)
	fmt.Println(segment, XYT)
	s2 := decompressXYT(2214, 11)
	fmt.Println(s2)
}

func TestXYT4DiffErr(t *testing.T) {
	//unix转换成时间戳
	t2 := time.Unix(1211075830, 0).UTC().Format(layout)
	fmt.Println(t2)
	//转二进制
	s2 := strconv.FormatInt(166, 2)
	s3 := strconv.FormatInt(2214&(1<<11-1), 2)
	fmt.Println(s2, s3)
	s := mttypes.SecondIndex{
		ID:      "1",
		StartTs: "1211075830",
		Segment: "32",
	}

	i, err := compressXYT(s.StartTs, 11)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}
	fmt.Println(i)

}

func TestEasy(t *testing.T) {
	a := "10100110"
	fmt.Println(11 - len(a))
	length := 11 - len(a)
	for i := 0; i < length; i++ {
		a = "0" + a
	}
	fmt.Println(a)
}
