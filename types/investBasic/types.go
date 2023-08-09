package investBasic

import (
	"fmt"
)

type EDirTrade uint

const (
	EMPTY EDirTrade = iota
	LONG
	SHORT
)

type ExchangeCD string

const (
	SSE   ExchangeCD = "SSE"
	SZSE  ExchangeCD = "SZSE"
	SHFE  ExchangeCD = "SHFE"
	CFFEX ExchangeCD = "CFFEX"
	INE   ExchangeCD = "INE"
	DCE   ExchangeCD = "DCE"
	CZCE  ExchangeCD = "CZCE"
	HKEX  ExchangeCD = "HKEX"
)

type TypeSecurity string

const (
	STOCK  TypeSecurity = "STOCK"
	FUND   TypeSecurity = "FUND"
	BOND   TypeSecurity = "BOND"
	FUTURE TypeSecurity = "FUTURE"
	IDX    TypeSecurity = "IDX"
)

type SSecInfo struct {
	Code       string       `json:"code"`
	ExchangeCD ExchangeCD   `json:"exchangecd"`
	Name       string       `json:"name"`
	Type       TypeSecurity `json:"type"`
}

func (s *SSecInfo) UniqueCode() string {
	return fmt.Sprintf("%s.%s", s.Code, s.ExchangeCD)
}

// SSecInfoChange 证券基本信息变动情况
type SSecInfoChange struct {
	// 变动时间
	Date string `json:"date"`
	// 变动后信息
	ExchangeCD ExchangeCD   `json:"exchangecd"`
	Type       TypeSecurity `json:"type"`
	Code       string       `json:"code"`
	Name       string       `json:"name"`
}

func (secInfoChange SSecInfoChange) InitFromSecInfo(secInfo SSecInfo, date string) *SSecInfoChange {
	secInfoChange.Date = date
	secInfoChange.Code = secInfo.Code
	secInfoChange.ExchangeCD = secInfo.ExchangeCD
	secInfoChange.Name = secInfo.Name
	secInfoChange.Type = secInfo.Type
	return &secInfoChange
}

// SMDKindleDailyFuture 日行情-期货
type SMDKindleDailyFuture struct {
	// 日期
	Date string `json:"date" bson:"date"`
	// 代码
	Code string `json:"code" bson:"code"`
	// 开
	Open float64 `json:"open" bson:"open"`
	// 高
	High float64 `json:"high" bson:"high"`
	// 低
	Low float64 `json:"low" bson:"low"`
	// 收
	Close float64 `json:"close" bson:"close"`
	// 结算
	Settle float64 `json:"settle" bson:"settle"`
	// 前结算
	PreSettle float64 `json:"preSettle" bson:"preSettle"`
	// 成交额
	Val float64 `json:"val" bson:"val"`
	// 成交量
	Vol int64 `json:"vol" bson:"vol"`
	// 持仓量
	NumHold int64 `json:"numHold" bson:"numHold"`
}

// SMDKindleDailyStock 日行情-股票
type SMDKindleDailyStock struct {
	// 日期
	Date string `json:"date" bson:"date"`
	// 代码
	Code string `json:"code" bson:"code"`
	// 开
	Open int64 `json:"open" bson:"open"`
	// 高
	High int64 `json:"high" bson:"high"`
	// 低
	Low int64 `json:"low" bson:"low"`
	// 收
	Close int64 `json:"close" bson:"close"`
	// 前收
	PreClose int64 `json:"preClose" bson:"preClose"`
	// 成交额
	Val float64 `json:"val" bson:"val"`
	// 成交量
	Vol int64 `json:"vol" bson:"vol"`
}

func (s SMDKindleDailyStock) String() string {
	return fmt.Sprintf("%s,%s,%d,%d,%d,%d,%d,%.2f,%d", s.Date, s.Code, s.Open, s.High, s.Low, s.Close, s.PreClose, s.Val, s.Vol)
}

// SETFNumHold ETF规模
type SETFNumHold struct {
	// 日期
	Date string
	// 交易所
	ExchangeCD string
	// 代码
	Code string
	//总份额
	NumHold float64
}

// SMktConnFlow 互联互通数据
type SMktConnFlow struct {
	// 日期
	Date string `json:"date" bson:"date"`
	// 方向：hk2sh,hk2sz
	Dir string `json:"dir" bson:"dir"`
	// 买入成交额
	ValBuy float64 `json:"valBuy" bson:"valBuy"`
	// 卖出成交额
	ValSell float64 `json:"valSell" bson:"valSell"`
	// 净流入
	ValNet float64 `json:"valNet" bson:"valNet"`
	// 累计净流入
	ValCum float64 `json:"valCum" bson:"valCum"`
	//净买入（委托）
	ValIn float64 `json:"valIn" bson:"valIn"`
}

type SFundNetValEM struct {
	// 代码
	Code string `json:"code" bson:"code"`
	// 日期
	FSRQ string `json:"FSRQ" bson:"FSRQ"`
	// 单位净值(普通基金)/万份收益（货基）
	DWJZ string `json:"DWJZ" bson:"DWJZ"`
	// 累计净值（普通基金）/7日年化收益率（伙计）
	LJJZ string `json:"LJJZ" bson:"LJJZ"`
	// 分红送配说明
	FHSP string `json:"FHSP" bson:"FHSP"`
	// 分红分拆值
	FHFCZ string `json:"FHFCZ" bson:"FHFCZ"`
	// 净值增长率
	JZZZL string `json:"JZZZL" bson:"JZZZL"`
}

type SMDTick struct {
	UniqueCode          string  `json:"uniqueCode"`
	DayTrade            string  `json:"dayTrade"`
	Timestamp           int64   `json:"timestamp"`
	PricePreClose       int64   `json:"pricePreClose"`
	PriceLimitLower     int64   `json:"priceLimitLower"`
	PriceLimitUpper     int64   `json:"priceLimitUpper"`
	PriceOpen           int64   `json:"priceOpen"`
	PriceHigh           int64   `json:"priceHigh"`
	PriceLow            int64   `json:"priceLow"`
	PriceLatest         int64   `json:"priceLatest"`
	Num                 int64   `json:"num"`
	Vol                 int64   `json:"vol"`
	Val                 float64 `json:"val"`
	NegVal              float64 `json:"negVal"`
	AskVols             []int64 `json:"askVols"`
	BidVols             []int64 `json:"bidVols"`
	AskPrices           []int64 `json:"askPrices"`
	BidPrices           []int64 `json:"bidPrices"`
	WeightedAvgAskPrice int64   `json:"weightedAvgAskPrice"`
	WeightedAvgBskPrice int64   `json:"weightedAvgBskPrice"`
	AskVolTotal         int64   `json:"askVolTotal"`
	BidVolTotal         int64   `json:"bidVolTotal"`
}

// SVolDDZR 成交量指标(DDZR)
type SVolDDZR struct {
	UcProduct string
	// 日期
	Date string
	// 成交量1
	Vol1 int64
	// 成交量2
	Vol2 int64
	// 成交量3
	Vol3 int64
}

// STradeDay 交易日信息
type STradeDay struct {
	Date       string `json:"date"`
	IsTradeDay bool   `json:"isTradeDay"`
	DateLast   string `json:"dateLast"`
	DateNext   string `json:"dateNext"`
}

// SSecCodeCCASS 互联互通港交所代码和a股对应关系
type SSecCodeCCASS struct {
	// 数据日期
	Date string `json:"date"`
	// 对应证券
	UniqueCode string `json:"uniqueCode"`
	// ccasscode
	CcassCode string `json:"ccassCode"`
	// 是否属于互联互通可买入
	Buyable bool `json:"buyable"`
}

// SHoldHKEX 港交所持股汇总
type SHoldHKEX struct {
	// 持股日期
	Date string `json:"date"`
	// 持股代码
	CodeCCASS string `json:"codeCCASS"`
	// 持股量
	HoldNum int64 `json:"holdNum"`
	// 持股占比
	HoldRate float64 `json:"holdRate"`
}

// SInvestorHKEX 港交所投资者信息
type SInvestorHKEX struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type SHoldDetailHKEX struct {
	// 持股日期
	Date string `json:"date"`
	// 持股代码
	CodeCCASS string `json:"codeCCASS"`
	// 投资者
	Investor string `json:"investor"`
	// 持股量
	HoldNum int64 `json:"holdNum"`
	// 持股占比
	HoldRate float64 `json:"holdRate"`
}
