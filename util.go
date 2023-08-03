package invest

import (
	"fmt"
	"github.com/rz1998/invest-basic/types/investBasic"
	"regexp"
	"sort"
	"strings"
	"time"
)

// GetCodeProduct 获取品种代码， 股票代码不变，期货只取英文部分
func GetCodeProduct(code string) string {
	regFuture := regexp.MustCompile(`^[A-Za-z]+`)
	codeProduct := regFuture.FindString(code)
	if codeProduct != "" {
		return codeProduct
	} else {
		return strings.Split(code, ".")[0]
	}
}

// GetSecInfo 唯一码转换为代码和交易所
func GetSecInfo(uniqueCode string) (code string, exchangeCD investBasic.ExchangeCD) {
	contents := strings.Split(uniqueCode, ".")
	if len(contents) > 1 {
		return contents[0], investBasic.ExchangeCD(contents[1])
	} else {
		return contents[0], ""
	}
}

// GetSecType 解析证券类型
func GetSecType(uniqueCode string) investBasic.TypeSecurity {
	code, exchangeCD := GetSecInfo(uniqueCode)
	switch exchangeCD {
	case investBasic.SHFE, investBasic.CFFEX, investBasic.INE, investBasic.DCE, investBasic.CZCE:
		return investBasic.FUTURE
	case investBasic.SSE:
		if code[:1] == "6" {
			return investBasic.STOCK
		} else if code[:1] == "5" {
			return investBasic.FUND
		} else if code[:3] == "000" {
			return investBasic.IDX
		} else {
			return investBasic.BOND
		}
	case investBasic.SZSE:
		if code[:3] == "399" {
			return investBasic.IDX
		} else if code[:1] == "0" || code[:1] == "2" || code[:1] == "3" {
			return investBasic.STOCK
		} else if code[:2] == "15" {
			return investBasic.FUND
		} else {
			return investBasic.BOND
		}
	default:
		return investBasic.STOCK
	}
}

// fromCCASS2Std 互联互通的港交所ccass代码转换为标准代码
func fromCCASS2Std(codeCCASS string) string {
	code := ""
	if len(codeCCASS) == 0 {
		fmt.Printf("FromCCASS2Std stopped by no ccass\n")
		return code
	}
	if codeCCASS[:2] == "77" {
		code = "300" + codeCCASS[2:]
	} else if codeCCASS[:2] == "78" {
		code = "301" + codeCCASS[2:]
	} else if codeCCASS[:1] == "7" {
		code = "00" + codeCCASS[1:]
	} else if codeCCASS[:1] == "9" {
		code = "60" + codeCCASS[1:]
	} else if codeCCASS[:2] == "30" {
		code = "688" + codeCCASS[2:]
	} else {
		fmt.Printf("FromCCASS2Std stopped by unhandled ccass %s\n", codeCCASS)
		return code
	}
	var exchangeCD investBasic.ExchangeCD
	if len(code) == 0 {
		return code
	} else if code[:1] == "6" {
		exchangeCD = investBasic.SSE
	} else {
		exchangeCD = investBasic.SZSE
	}
	return fmt.Sprintf("%s.%s", code, exchangeCD)
}

func FromCCASS2Std(mapCCASS map[string]string, codeCCASS string) string {
	if uniqueCode, hasCCASS := mapCCASS[codeCCASS]; hasCCASS {
		return uniqueCode
	} else {
		return fromCCASS2Std(codeCCASS)
	}
}

// IsCodeFutureProduct 判断期货代码是否是品种
/*
 * 是： 品种
 * 否： 合约
 */
func IsCodeFutureProduct(code string) bool {
	reg := regexp.MustCompile(`\w+\d+`)
	if reg.Match([]byte(code)) {
		return false
	} else {
		return true
	}
}

// MergeContract2IdxMD 合约行情合并成指数行情 uniqueCode合约唯一码
func MergeContract2IdxMD(uniqueCode string, mdCons [][]*investBasic.SMDTick) []*investBasic.SMDTick {
	var results []*investBasic.SMDTick
	// 合并排序
	if mdCons == nil {
		return results
	}
	// 合并之后排序
	var mds []*investBasic.SMDTick
	for _, mdCon := range mdCons {
		mds = append(mds, mdCon...)
	}
	sort.Slice(mds, func(i, j int) bool {
		return mds[i].Timestamp < mds[j].Timestamp
	})
	//
	nameFieldSums := []string{"Val", "NegVal"}
	mapMD := make(map[string]*investBasic.SMDTick)
	var volCache int64
	for _, md := range mds {
		// 更新行情
		mapMD[md.UniqueCode] = md

		// 数据充足才开始计算
		if len(mapMD) < len(mdCons) {
			volCache += md.Vol
			continue
		}

		mdMembers := make([]interface{}, len(mapMD))
		i := 0
		for _, v := range mapMD {
			mdMembers[i] = v
			i++
		}
		mapIdx := SumStructs(nameFieldSums, mdMembers)
		priceWAVG := WAvgStruct("NegVal", "PriceLatest", mdMembers)
		// 处理初始成交量
		vol := md.Vol
		if volCache > 0 {
			vol += volCache
			volCache = 0
		}
		mdIdx := investBasic.SMDTick{
			UniqueCode:  uniqueCode,
			DayTrade:    md.DayTrade,
			Timestamp:   md.Timestamp,
			PriceLatest: int64(priceWAVG),
			Num:         md.Num,
			Vol:         vol,
			Val:         mapIdx["Val"],
			NegVal:      mapIdx["NegVal"],
		}
		results = append(results, &mdIdx)
	}
	return results
}

// MergeTick2Minute tick转换为分钟行情
func MergeTick2Minute(mdTicks []*investBasic.SMDTick) []*investBasic.SMDTick {
	var result []*investBasic.SMDTick
	var highDaily, lowDaily int64
	arrayTicks := make([]interface{}, len(mdTicks))
	for i, md := range mdTicks {
		arrayTicks[i] = md
	}
	arraySliced := SliceTimeSeries("Timestamp", 1*time.Minute, arrayTicks)
	nameFieldSums := []string{"Vol", "NegVal"}
	result = make([]*investBasic.SMDTick, len(arraySliced))
	for i, array := range arraySliced {
		md, ok := MergeTimeSeries(nameFieldSums, array).(*investBasic.SMDTick)
		if ok {
			result[i] = md
			md.PriceOpen = array[0].(*investBasic.SMDTick).PriceLatest
			//最高最低价
			var priceHigh, priceLow, priceLatest, minuteHigh, minuteLow int64
			for _, data := range array {
				md0 := data.(*investBasic.SMDTick)
				priceHigh = md0.PriceHigh
				priceLow = md0.PriceLow
				priceLatest = md0.PriceLatest
				// 初始化
				if highDaily == 0 {
					highDaily = priceHigh
				}
				if lowDaily == 0 {
					lowDaily = priceLow
				}
				// 最高最低有变化时，取最高最低
				if highDaily < priceHigh {
					highDaily = priceHigh
					minuteHigh = priceHigh
				}
				if lowDaily > priceLow {
					lowDaily = priceLow
					minuteLow = priceLow
				}
				// 分钟内最新价更新
				if minuteHigh == 0 || minuteHigh < priceLatest {
					minuteHigh = priceLatest
				}
				if minuteLow == 0 || minuteLow > priceLatest {
					minuteLow = priceLatest
				}
			}
			md.PriceHigh = minuteHigh
			md.PriceLow = minuteLow
		}
	}
	return result
}
