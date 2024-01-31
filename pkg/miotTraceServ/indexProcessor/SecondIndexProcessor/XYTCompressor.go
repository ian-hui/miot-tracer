package secondindexprocessor

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

func splitXYT2StartEnd(combined int64) (start_ts int64, end_ts int64) {
	end_ts = combined & int64(mttypes.TS_LEN)
	start_ts = combined >> 16
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

// 合并segment和XYT
func splitSegmentAndXYT(combined int64) (segment int64, XYT int64) {
	segment = combined & int64(mttypes.SEGMENT_LEN)
	XYT = combined >> 8
	return
}

// 分离nextnode和segment和XYT
func splitAll(combined int64) (next_node int64, segment int64, XYT int64) {
	next_node = combined & int64(mttypes.NEXT_NODE_LEN)
	xyt_segment := combined >> 8
	segment = xyt_segment & int64(mttypes.SEGMENT_LEN)
	XYT = xyt_segment >> 8
	return
}

// 合并XYT和segment
func combineXYTAndSegment(XYT int64, segment int64) int64 {
	return XYT<<8 | segment
}

// 合并整个secondindex
func combineAll(XYT int64, segment int64, next_node int64) int64 {
	xyt_segment := XYT<<8 | segment
	return xyt_segment<<8 | next_node
}

// 解压secondindex
func decompressSecondIndex(combined int64) (start_ts string, end_ts string, segment string, next_node string) {
	unprocessed_next_node, unprocessed_segment, unprocessed_XYT := splitAll(combined)
	segment = strconv.FormatInt(unprocessed_segment, 10)
	next_node = strconv.FormatInt(unprocessed_next_node, 10)
	strat, end := splitXYT2StartEnd(unprocessed_XYT)
	start_ts = decompressXYT(strat)
	end_ts = decompressXYT(end)
	return
}

func String2UnixTimestamp(ts string) (int64, error) {
	i, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		iotlog.Errorln("strconv.ParseInt failed, err:", err)
		return 0, err
	}
	return i, nil
}

func combineXYT(start_ts int64, end_ts int64) int64 {
	return start_ts<<2 | end_ts
}
