package impl

import (
	"encoding/json"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const (
	MEMB  = "B"
	MEMKB = "KB"
	MEMMB = "MB"
	MEMGB = "GB"
	MEMTB = "TB"
	MEMPB = "PB"
)
const (
	SIZEB = 1 << (iota * 10)
	SIZEKB
	SIZEMB
	SIZEGB
	SIZETB
	SIZEPB
)

func parseSize(s string) (int64, string) {
	//拆分正则表达
	re, _ := regexp.Compile("[0-9]+")
	//获取单位
	unit := string(re.ReplaceAll([]byte(s), []byte("")))
	num, _ := strconv.ParseInt(strings.Replace(s, unit, "", 1), 10, 64)
	unit = strings.ToUpper(unit)
	switch unit {
	case MEMB:
		num = num * SIZEB
	case MEMKB:
		num = num * SIZEKB
	case MEMMB:
		num = num * SIZEMB
	case MEMGB:
		num = num * SIZEGB
	case MEMTB:
		num = num * SIZETB
	case MEMPB:
		num = num * SIZEPB
	default:
		num = 0
	}
	if num == 0 {
		log.Printf("暂不支持size为%v\n,", s)
		num = 100 * SIZEMB
		unit = MEMMB
		s = "100MB"
	}

	return num, s
}

func getValSize(val interface{}) int64 {
	bytes, _ := json.Marshal(val)
	//size := unsafe.Sizeof(val)
	return int64(len(bytes))
}
