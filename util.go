package invest

import (
	"fmt"
	"github.com/rz1998/invest-basic/types/investBasic"
	"regexp"
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
