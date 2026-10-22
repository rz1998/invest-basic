/*
 *                        ::
 *                       :;J7, :,                        ::;7:
 *                       ,ivYi, ,                       ;LLLFS:
 *                       :iv7Yi                       :7ri;j5PL
 *                      ,:ivYLvr                    ,ivrrirrY2X,
 *                      :;r@Wwz.7r:                :ivu@kexianli.
 *                     :iL7::,:::iiirii:ii;::::,,irvF7rvvLujL7ur
 *                    ri::,:,::i:iiiiiii:i:irrv177JX7rYXqZEkvv17
 *                 ;i:, , ::::iirrririi:i:::iiir2XXvii;L8OGJr71i
 *               :,, ,,:   ,::ir@mingyi.irii:i:::j1jri7ZBOS7ivv,
 *                  ,::,    ::rv77iiiriii:iii:i::,rvLq@huhao.Li
 *              ,,      ,, ,:ir7ir::,:::i;ir:::i:i::rSGGYri712:
 *            :::  ,v7r:: ::rrv77:, ,, ,:i7rrii:::::, ir7ri7Lri
 *           ,     2OBBOi,iiir;r::        ,irriiii::,, ,iv7Luur:
 *         ,,     i78MBBi,:,:::,:,  :7FSL: ,iriii:::i::,,:rLqXv::
 *         :      iuMMP: :,:::,:ii;2GY7OBB0viiii:i:iii:i:::iJqL;::
 *        ,     ::::i   ,,,,, ::LuBBu BBBBBErii:i:i:i:i:i:i:r77ii
 *       ,       :       , ,,:::rruBZ1MBBqi, :,,,:::,::::::iiriri:
 *      ,               ,,,,::::i:  @arqiao.       ,:,, ,:::ii;i7:
 *     :,       rjujLYLi   ,,:::::,:::::::::,,   ,:i,:,,,,,::i:iii
 *     ::      BBBBBBBBB0,    ,,::: , ,:::::: ,      ,,,, ,,:::::::
 *     i,  ,  ,8BMMBBBBBBi     ,,:,,     ,,, , ,   , , , :,::ii::i::
 *     :      iZMOMOMBBM2::::::::::,,,,     ,,,,,,:,,,::::i:irr:i:::,
 *     i   ,,:;u0MBMOG1L:::i::::::  ,,,::,   ,,, ::::::i:i:iirii:i:i:
 *     :    ,iuUuuXUkFu7i:iii:i:::, :,:,: ::::::::i:i:::::iirr7iiri::
 *     :     :rk@Yizero.i:::::, ,:ii:::::::i:::::i::,::::iirrriiiri::,
 *      :      5BMBBBBBBSr:,::rv2kuii:::iii::,:i:,, , ,,:,:i@petermu.,
 *           , :r50EZ8MBBBBGOBBBZP7::::i::,:::::,: :,:,::i;rrririiii::
 *               :jujYY7LS0ujJL7r::,::i::,::::::::::::::iirirrrrrrr:ii:
 *            ,:  :@kevensun.:,:,,,::::i:i:::::,,::::::iir;ii;7v77;ii;i,
 *            ,,,     ,,:,::::::i:iiiii:i::::,, ::::iiiir@xingjief.r;7:i,
 *         , , ,,,:,,::::::::iiiiiiiiii:,:,:::::::::iiir;ri7vL77rrirri::
 *          :,, , ::::::::i:::i:::i:i::,,,,,:,::i:i:::iir;@Secbone.ii:::
 *
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2024-05-16 02:16:45
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2024-06-01 01:21:02
 * @FilePath: /stra/home/rz1998/workspace/github.com/rz1998/invest/investBasic/utilFuture.go
 * @Description:
 *
 */

package investBasic

import (
	"fmt"
	"math"
	"regexp"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

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

/*
	GetDateDelivery 计算交割日(指定年份)

param
  - year 哪一年
  - weekOfMonth 周一起的第几周
  - dayOfWeek Sunday起
  - monthOfYear 0起

return
  - 指定判断的日期

example
  - GetDateDelivery(yearNow, 3, int(time.Friday), int(time.August)) 返回今年8月第3个星期五的日期
*/
func GetDateDelivery(year, weekOfMonth, dayOfWeek, monthOfYear int) time.Time {
	strMethod := "GetDateDelivery"
	// 获取指定月第一天日期
	firstDayOfMonth, err := time.Parse("2006-1-02", fmt.Sprintf("%d-%d-01", year, monthOfYear))
	if err != nil {
		logx.Errorf("%s stopped by %v\n", strMethod, err)
		return firstDayOfMonth
	}
	// 根据第一天的星期几，计算目标天数差距
	diffWeekday := time.Weekday(dayOfWeek) - firstDayOfMonth.Weekday()
	diffDays := weekOfMonth*7 - int(math.Ceil(float64(diffWeekday)/7))*7 + int(diffWeekday)
	return firstDayOfMonth.AddDate(0, 0, diffDays)
}

// FromDeliveryMonthFuture 用交割月份生成合约代码
/*
 * @param ucProduct 品种唯一码
 * @param year 交割年
 * @param month 交割月
 * @return uniqueCode
 */
func FromDeliveryMonthFuture(ucProduct string, year, month int) string {
	strMethod := "FromDeliveryMonthFuture"
	var uniqueCode string
	codeProduct, exchange := GetSecInfo(ucProduct)
	if len(codeProduct) == 0 || len(exchange) == 0 {
		logx.Errorf("%s stopped by %s %s\n", strMethod, "error ucProduct", ucProduct)
		return uniqueCode
	}
	date, err := time.Parse("2006-1-02", fmt.Sprintf("%d-%d-01", year, month))
	if err != nil {
		logx.Errorf("%s stopped by %v\n", strMethod, err)
		return uniqueCode
	}
	switch exchange {
	// 正常交易所
	case CFFEX, SHFE, DCE, INE:
		uniqueCode = fmt.Sprintf("%s%s.%s", codeProduct, date.Format("0601"), exchange)
	// 郑商所装怪
	case CZCE:
		uniqueCode = fmt.Sprintf("%s%s.%s", codeProduct, date.Format("0601")[1:], exchange)
	default:
		logx.Errorf("%s stopped by unhandled exchangeCD %s\n", strMethod, ucProduct)
	}
	return uniqueCode
}

// 主力合约判断
func QryMainContract(product string, date time.Time) string {
	// 返回交割日
	dateDelivery := GetDateDelivery(date.Year(), 3, int(time.Friday), int(date.Month()))
	dateContract := date.AddDate(0, 0, -date.Day()+1)
	limit := date.AddDate(0, 0, 2)
	var tarDate time.Time
	if limit.Before(dateDelivery) {
		// 当月合约到期日 - 当前交易日 > 2: 当月合约
		tarDate = dateContract
	} else {
		// 当月合约到期日 - 当前交易日  2: 下月合约
		tarDate = dateContract.AddDate(0, 1, 0)
	}
	return FromDeliveryMonthFuture(fmt.Sprintf("%s.%s", product, CFFEX), tarDate.Year(), int(tarDate.Month()))
}
