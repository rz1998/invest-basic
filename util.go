/*
 *                        _oo0oo_
 *                       o8888888o
 *                       88" . "88
 *                       (| -_- |)
 *                       0\  =  /0
 *                     ___/`---'\___
 *                   .' \\|     |// '.
 *                  / \\|||  :  |||// \
 *                 / _||||| -:- |||||- \
 *                |   | \\\  - /// |   |
 *                | \_|  ''\---/''  |_/ |
 *                \  .-\__  '-'  ___/-. /
 *              ___'. .'  /--.--\  `. .'___
 *           ."" '<  `.___\_<|>_/___.' >' "".
 *          | | :  `- \`.;`\ _ /`;.`/ - ` : | |
 *          \  \ `_.   \_ __\ /__ _/   .-` /  /
 *      =====`-.____`.___ \_____/___.-`___.-'=====
 *                        `=---='
 *
 *
 *      ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 *
 *            佛祖保佑     永不宕机     永无BUG
 *
 *        佛曰:
 *                写字楼里写字间，写字间里程序员；
 *                程序人员写程序，又拿程序换酒钱。
 *                酒醒只在网上坐，酒醉还来网下眠；
 *                酒醉酒醒日复日，网上网下年复年。
 *                但愿老死电脑间，不愿鞠躬老板前；
 *                奔驰宝马贵者趣，公交自行程序员。
 *                别人笑我忒疯癫，我笑自己命太贱；
 *                不见满街漂亮妹，哪个归得程序员？
 *
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2024-04-08 08:40:11
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2024-08-15 20:26:55
 * @FilePath: /stra/home/rz1998/workspace/github.com/rz1998/invest/investBasic/util.go
 * @Description:
 *
 */

package investBasic

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"
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
//
// 后缀接受两种写法, 归一到 Go 交易所码 (SSE/SZSE/BSE):
//
//	Go 形式   : 600000.SSE  000001.SZSE  430047.BSE
//	QMT 形式  : 600000.SH   000001.SZ    430047.BJ
//
// 为什么要认 QMT 形式 (2026-10-22, best-shm 后缀统一): best 策略在信号入口
// 把代码转成 QMT 标准后缀后, 全链路 (持仓 map / 行情 / 报单) 都用 QMT 形式。
// 本函数是共享库里所有按交易所分支的入口 (GetMinVol / GetTickVol /
// GetTickPrice / GetVolSellableFromAcPos / GetSecType 都先调它)。不认
// .SH/.SZ/.BJ 的话, 这些分支全掉进 default: 最小下单量算成 1、tick 算成 1、
// 可卖量公式走错支 -- 全是静默算错。
//
// 大写归一: 后缀大小写不敏感 (QMT 回报告诉我们的是 "SH", 不是 "sh")。
func GetSecInfo(uniqueCode string) (code string, exchangeCD ExchangeCD) {
	contents := strings.Split(uniqueCode, ".")
	if len(contents) > 1 {
		suffix := strings.ToUpper(contents[1])
		switch suffix {
		case "SH":
			return contents[0], SSE
		case "SZ":
			return contents[0], SZSE
		case "BJ":
			return contents[0], BSE
		default:
			return contents[0], ExchangeCD(contents[1])
		}
	}
	return contents[0], ""
}

// GetSecType 解析证券类型
func GetSecType(uniqueCode string) TypeSecurity {
	code, exchangeCD := GetSecInfo(uniqueCode)
	switch exchangeCD {
	case SHFE, CFFEX, INE, DCE, CZCE:
		return FUTURE
	case SSE:
		if code[:1] == "6" {
			return STOCK
		} else if code[:1] == "5" {
			return FUND
		} else if code[:3] == "000" {
			return IDX
		} else {
			return BOND
		}
	case SZSE:
		if code[:3] == "399" {
			return IDX
		} else if code[:1] == "0" || code[:1] == "2" || code[:1] == "3" {
			return STOCK
		} else if code[:2] == "15" {
			return FUND
		} else {
			return BOND
		}
	default:
		return STOCK
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
	var exchangeCD ExchangeCD
	if len(code) == 0 {
		return code
	} else if code[:1] == "6" {
		exchangeCD = SSE
	} else {
		exchangeCD = SZSE
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

// 把小数精度转换为整数
func FromFloat(idx float64, accuracy int) int64 {
	return int64(math.Round(idx * math.Pow10(accuracy)))
}

// 计算指标中位数
func CalMedian[T int | int32 | int64 | float32 | float64](data []T) float64 {
	sort.Slice(data, func(i, j int) bool {
		return data[i] < data[j]
	})
	if len(data)%2 == 0 {
		// 偶数取中间两个数的均值
		return (float64(data[len(data)/2-1]) + float64(data[len(data)/2])) / 2
	} else {
		// 奇数取中间的数
		return float64(data[len(data)/2])
	}
}

func FromMD2Kindle(md *SMDTick) *SMDKindleDailyStock {
	if md == nil {
		return nil
	}
	return &SMDKindleDailyStock{
		Date:     md.DayTrade,
		Code:     md.UniqueCode,
		Open:     md.PriceOpen,
		High:     md.PriceHigh,
		Low:      md.PriceLow,
		Close:    md.PriceLatest,
		PreClose: md.PricePreClose,
		Vol:      md.Vol,
		Val:      md.Val,
	}
}
