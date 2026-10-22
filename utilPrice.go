/*
 *  ┌───┐   ┌───┬───┬───┬───┐ ┌───┬───┬───┬───┐ ┌───┬───┬───┬───┐ ┌───┬───┬───┐
 *  │Esc│   │ F1│ F2│ F3│ F4│ │ F5│ F6│ F7│ F8│ │ F9│F10│F11│F12│ │P/S│S L│P/B│  ┌┐    ┌┐    ┌┐
 *  └───┘   └───┴───┴───┴───┘ └───┴───┴───┴───┘ └───┴───┴───┴───┘ └───┴───┴───┘  └┘    └┘    └┘
 *  ┌───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───────┐ ┌───┬───┬───┐ ┌───┬───┬───┬───┐
 *  │~ `│! 1│@ 2│# 3│$ 4│% 5│^ 6│& 7│* 8│( 9│) 0│_ -│+ =│ BacSp │ │Ins│Hom│PUp│ │N L│ / │ * │ - │
 *  ├───┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─────┤ ├───┼───┼───┤ ├───┼───┼───┼───┤
 *  │ Tab │ Q │ W │ E │ R │ T │ Y │ U │ I │ O │ P │{ [│} ]│ | \ │ │Del│End│PDn│ │ 7 │ 8 │ 9 │   │
 *  ├─────┴┬──┴┬──┴┬──┴┬──┴┬──┴┬──┴┬──┴┬──┴┬──┴┬──┴┬──┴┬──┴─────┤ └───┴───┴───┘ ├───┼───┼───┤ + │
 *  │ Caps │ A │ S │ D │ F │ G │ H │ J │ K │ L │: ;│" '│ Enter  │               │ 4 │ 5 │ 6 │   │
 *  ├──────┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴─┬─┴────────┤     ┌───┐     ├───┼───┼───┼───┤
 *  │ Shift  │ Z │ X │ C │ V │ B │ N │ M │< ,│> .│? /│  Shift   │     │ ↑ │     │ 1 │ 2 │ 3 │   │
 *  ├─────┬──┴─┬─┴──┬┴───┴───┴───┴───┴───┴──┬┴───┼───┴┬────┬────┤ ┌───┼───┼───┐ ├───┴───┼───┤ E││
 *  │ Ctrl│    │Alt │         Space         │ Alt│    │    │Ctrl│ │ ← │ ↓ │ → │ │   0   │ . │←─┘│
 *  └─────┴────┴────┴───────────────────────┴────┴────┴────┴────┘ └───┴───┴───┘ └───────┴───┴───┘
 *
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2024-05-14 02:45:40
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2024-09-20 10:35:11
 * @FilePath: /tradeNew/home/rz1998/workspace/github.com/rz1998/invest/trade/tradeBasic/utilPrice.go
 * @Description: 交易过程中对价格的计算
 *
 */

package investBasic

import (
	"math"

	"github.com/zeromicro/go-zero/core/logx"
)

// 根据代码返回最小价格变动单位
func GetTickPrice(uniqueCode string) int64 {
	code, exchangeCD := GetSecInfo(uniqueCode)
	var tick int64
	switch exchangeCD {
	case SSE, SZSE, BSE:
		// 沪深市场
		str1 := code[:1]
		str2 := code[:2]
		if str1 == "0" || str1 == "3" || str1 == "6" ||
			str1 == "4" || str1 == "8" || str1 == "9" {
			// 股票，100一跳
			tick = 100
		} else if exchangeCD == SSE && str2 == "11" {
			// 上交所可转债，100一跳
			tick = 100
		} else {
			tick = 10
		}
	case CFFEX:
		if code[:1] == "I" {
			// 股指0.2一跳
			tick = 2000
		} else if code[:1] == "T" {
			// 国债0.005一跳
			tick = 50
		}
	default:
		logx.Errorf("CalTarPrice unhandled exchangeCD %v\n", exchangeCD)
	}
	return tick
}

// CalLimitPrice 计算目标涨跌幅价(股票四舍五入)
func CalLimitPrice(rateStop float64, preClose, tick int64) int64 {
	if tick <= 0 {
		return 0
	}
	return int64(math.Round(float64(preClose)/float64(tick)*(1.0+rateStop))) * tick
}

/*
CalTarPrice 根据证券代码计算目标涨跌幅度价格

rateTar 限制幅度，比如10%则填写0.1

priceLimit 涨跌停价，rateTar>0用涨停价，rateTar<0用跌停价
*/
func CalTarPrice(uniqueCode string, rateTar float64, price, priceLimit int64) int64 {
	tick := GetTickPrice(uniqueCode)
	if priceLimit < 0 {
		// 忽略限制
		return CalLimitPrice(rateTar, price, tick)
	} else {
		if rateTar == 0 {
			return price
		} else if rateTar > 0 {
			if price > priceLimit {
				logx.Infof("CalTarPrice error, price over limitUpper! rataTar %f price %d priceLimit %d\n", rateTar, price, priceLimit)
				return price
			} else {
				return int64(math.Min(float64(priceLimit), float64(CalLimitPrice(rateTar, price, tick))))
			}
		} else {
			if price < priceLimit {
				logx.Infof("CalTarPrice error, price under limitLower! rataTar %f price %d priceLimit %d\n", rateTar, price, priceLimit)
				return price
			} else {
				return int64(math.Max(float64(priceLimit), float64(CalLimitPrice(rateTar, price, tick))))
			}
		}
	}
}
