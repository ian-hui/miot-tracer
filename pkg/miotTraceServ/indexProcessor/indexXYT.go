package indexprocessor

import (
	mttypes "miot_tracing_go/mtTypes"
	"strconv"
)

// 压缩XYT(helper function)
func compressXYT(ts string) int64 {
	//binnum
	binNum := genBinNum(ts)
	//elementCode （11位）
	elementCode := genElementCode(binNum, ts, 11)
	//把binNum和elementCode合并
	combined := binNum<<11 | elementCode
	return combined
}

// 解压XYT
func decompressXYT(combined int64) (unix_ts string) {
	//先取出binNum
	binNum := combined >> 11
	//还原时间戳
	binnum_start_ts := mttypes.REF_TIME.Unix() + binNum*mttypes.BIN_LEN
	binnum_end_ts := binnum_start_ts + mttypes.BIN_LEN
	//取出elementCode
	elementCode := combined & int64(mttypes.ELEMENTCODE_LEN)
	//用二分的方法找到startTS的位置
	s := strconv.FormatInt(elementCode, 2)
	//补0
	for i := 0; i < 11-len(s); i++ {
		s = "0" + s
	}
	mid := (binnum_start_ts + binnum_end_ts) / 2
	for i := 0; i < len(s); i++ {
		if s[i] == '0' {
			mid = (binnum_start_ts + mid) / 2
		} else {
			mid = (mid + binnum_end_ts) / 2
		}
	}
	unix_ts = strconv.FormatInt(mid, 10)
	return
}

func genBinNum(ts string) int64 {
	//binnum
	i, err := String2UnixTimestamp(ts)
	if err != nil {
		iotlog.Errorln("String2UnixTimestamp failed, err:", err)
		return 0
	}
	diff := i - mttypes.REF_TIME.Unix()
	// 用差值除以 Bin Len 来计算 BinNum
	binNum := diff / mttypes.BIN_LEN
	//取低5位
	binNum = binNum & 0x1f
	return binNum
}

func genElementCode(binNum int64, startTS string, max_elementcode_len int) int64 {
	//还原时间戳
	binnum_start_ts := mttypes.REF_TIME.Unix() + binNum*mttypes.BIN_LEN
	binnum_end_ts := binnum_start_ts + mttypes.BIN_LEN
	ts, err := strconv.Atoi(startTS)
	if err != nil {
		iotlog.Errorln("strconv.Atoi failed, err:", err)
		return 0
	}
	ts64 := int64(ts)
	//用二分的方法找到startTS的位置
	s := ""
	mid := (binnum_start_ts + binnum_end_ts) / 2
	for {
		if mid == binnum_start_ts || len(s) == max_elementcode_len {
			break
		}
		if mid > ts64 {
			s += "0"
			mid = (binnum_start_ts + mid) / 2
		} else {
			s += "1"
			mid = (mid + binnum_end_ts) / 2
		}
	}
	//把s转成二进制
	bin, err := strconv.ParseInt(s, 2, 64)
	if err != nil {
		iotlog.Errorln("strconv.ParseInt failed, err:", err)
		return 0
	}
	//只要低11位
	bin = bin & int64(mttypes.ELEMENTCODE_LEN)
	return bin
}
