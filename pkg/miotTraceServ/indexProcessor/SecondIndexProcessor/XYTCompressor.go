package secondindexprocessor

import (
	"errors"
	mttypes "miot_tracing_go/mtTypes"
	"strconv"
)

//------------------------------------------------------------------
//------------------------compress and decompress-------------------
//------------------------------------------------------------------

func VariableLengthCompress(ts string, start_ts string) (int64, error) {
	// 转换成int64
	ts_int64, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		iotlog.Errorln("strconv.ParseInt failed, err:", err)
		return 0, err
	}
	first_index_start_ts_int64, err := strconv.ParseInt(start_ts, 10, 64)
	if err != nil {
		iotlog.Errorln("strconv.ParseInt failed, err:", err)
		return 0, err
	}
	// 计算差值
	diff := ts_int64 - first_index_start_ts_int64
	if diff <= 0 {
		iotlog.Errorln("diff <= 0")
		return 0, errors.New("diff <= 0")
	}
	if diff < int64(mttypes.VARIABLE_CHECK_LEN) {
		//采用秒索引
		exact_index, err := CompressExactIndex(diff)
		if err != nil {
			iotlog.Errorln("exactIndex failed, err:", err)
			return 0, err
		}
		return exact_index, nil
	} else {
		i, err := compressXYT(ts, 10)
		if err != nil {
			iotlog.Errorln("compressXYT failed, err:", err)
			return 0, err
		}
		return (1<<15 | i), nil
	}
}

func VariableLengthDecompress(combined int64, start_ts int64) (string, error) {
	mode_code := combined >> 15
	if mode_code == 0 {
		//采用秒索引
		exact_index, err := DecompressExactIndex(combined)
		if err != nil {
			iotlog.Errorln("DecompressExactIndex failed, err:", err)
			return "", err
		}
		exact_index += start_ts
		return strconv.FormatInt(exact_index, 10), nil
	} else {
		//采用XYT
		i := combined & int64(mttypes.VARIABLE_CHECK_LEN)
		return decompressXYT(i, 10), nil
	}
}

// 精确索引
func CompressExactIndex(diff int64) (int64, error) {
	// 用15bit保存差值
	return 0<<15 | diff, nil
}

func DecompressExactIndex(combined int64) (int64, error) {
	return combined & int64(mttypes.VARIABLE_CHECK_LEN), nil
}

// 压缩XYT(helper function)
func compressXYT(ts string, max_elementcode_len int) (int64, error) {
	//binnum（5位）
	binNum, err := genBinNum(ts)
	if err != nil {
		iotlog.Errorln("genBinNum failed, err:", err)
		return 0, err
	}
	//elementCode （max_elementcode_len位）
	elementCode := genElementCode(binNum, ts, max_elementcode_len)
	//把binNum和elementCode合并
	combined := binNum<<max_elementcode_len | elementCode
	return combined, nil
}

// 解压XYT
func decompressXYT(combined int64, max_elementcode_len int) (unix_ts string) {
	//先取出binNum
	binNum := combined >> max_elementcode_len
	//还原时间戳
	binnum_start_ts := mttypes.REF_TIME.UTC().Unix() + binNum*mttypes.BIN_LEN
	binnum_end_ts := binnum_start_ts + mttypes.BIN_LEN

	// fmt.Println("s", time.Unix(binnum_start_ts, 0).UTC(), time.Unix(binnum_end_ts, 0).UTC())
	//取出elementCode
	elementCode := combined & int64(1<<max_elementcode_len-1)
	//用二分的方法找到startTS的位置
	s := strconv.FormatInt(elementCode, 2)
	//补0
	for i := 0; i < max_elementcode_len-len(s); i++ {
		s = "0" + s
	}
	mid := (binnum_start_ts + binnum_end_ts) / 2
	for i := 0; i < len(s); i++ {
		if s[i] == '0' {
			binnum_end_ts = mid
			mid = (binnum_start_ts + mid) / 2
		} else {
			binnum_start_ts = mid
			mid = (mid + binnum_end_ts) / 2
		}
	}
	//如果最后一位是0，就取start_ts，否则取end_ts
	if s[len(s)-1] == '0' {
		unix_ts = strconv.FormatInt(binnum_start_ts, 10)
	} else {
		unix_ts = strconv.FormatInt(mid, 10)
	}
	return
}

// 解压secondindex
func decompressSecondIndex(combined int64) (start_ts string, end_ts string, segment string, next_node string, err error) {
	unprocessed_next_node, unprocessed_segment, unprocessed_XYT := splitAll(combined)

	segment = strconv.FormatInt(unprocessed_segment, 10)
	next_node = strconv.FormatInt(unprocessed_next_node, 10)
	start, end := splitXYT2StartEnd(unprocessed_XYT)

	start_ts = decompressXYT(start, 11)
	start_ts_64, err := strconv.ParseInt(start_ts, 10, 64)
	if err != nil {
		iotlog.Errorln("strconv.Atoi failed, err:", err)
		return
	}

	end_ts, err = VariableLengthDecompress(end, start_ts_64)
	if err != nil {
		iotlog.Errorln("VariableLengthDecompress failed, err:", err)
		return
	}

	return

}

//------------------------------------------------------------------
//------------------------binNum and elementCode--------------------
//------------------------------------------------------------------

func genBinNum(ts string) (int64, error) {
	//binnum
	i, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		iotlog.Errorln("strconv.ParseInt failed, err:", err)
		return 0, err
	}
	diff := i - mttypes.REF_TIME.UTC().Unix()
	// 用差值除以 Bin Len 来计算 BinNum
	binNum := diff / mttypes.BIN_LEN
	//取低5位
	binNum = binNum & 0x1f
	return binNum, nil
}

func genElementCode(binNum int64, startTS string, max_elementcode_len int) int64 {
	//还原时间戳
	binnum_start_ts := mttypes.REF_TIME.UTC().Unix() + binNum*mttypes.BIN_LEN
	binnum_end_ts := binnum_start_ts + mttypes.BIN_LEN

	ts, err := strconv.Atoi(startTS)
	if err != nil {
		iotlog.Errorln("strconv.Atoi failed, err:", err)
		return 0
	}
	ts64 := int64(ts)
	//用二分的方法找到startTS的位置
	s := ""
	mid := (binnum_end_ts-binnum_start_ts)/2 + binnum_start_ts
	for {
		if mid == binnum_start_ts || len(s) == max_elementcode_len {
			break
		}
		if mid > ts64 {
			s += "0"
			binnum_end_ts = mid
			mid = (binnum_start_ts + mid) / 2
		} else {
			s += "1"
			binnum_start_ts = mid
			mid = (mid + binnum_end_ts) / 2
		}
		// fmt.Println(time.Unix(mid, 0).UTC())
		// fmt.Println(s)

	}
	// fmt.Println(len(s))
	//把s转成二进制
	bin, err := strconv.ParseInt(s, 2, 64)
	if err != nil {
		iotlog.Errorln("strconv.ParseInt failed, err:", err)
		return 0
	}
	//只要低maxlen位
	bin = bin & int64((1<<max_elementcode_len)-1)
	return bin
}

//------------------------------------------------------------------
//------------------------split and  combine------------------------
//------------------------------------------------------------------

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

func splitXYT2StartEnd(combined int64) (start_ts int64, end_ts int64) {
	end_ts = combined & int64(mttypes.TS_LEN)
	start_ts = combined >> 16
	return
}

func combineStartXYTAndEndXYT(start_ts int64, end_ts int64) int64 {
	return start_ts<<16 | end_ts
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
