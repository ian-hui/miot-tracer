package secondindexprocessor

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

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
