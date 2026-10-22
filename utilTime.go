/*
 *                                                     __----~~~~~~~~~~~------___
 *                                    .  .   ~~//====......          __--~ ~~
 *                    -.            \_|//     |||\\  ~~~~~~::::... /~
 *                 ___-==_       _-~o~  \/    |||  \\            _/~~-
 *         __---~~~.==~||\=_    -_--~/_-~|-   |\\   \\        _/~
 *     _-~~     .=~    |  \\-_    '-~7  /-   /  ||    \      /
 *   .~       .~       |   \\ -_    /  /-   /   ||      \   /
 *  /  ____  /         |     \\ ~-_/  /|- _/   .||       \ /
 *  |~~    ~~|--~~~~--_ \     ~==-/   | \~--===~~        .\
 *           '         ~-|      /|    |-~\~~       __--~~
 *                       |-~~-_/ |    |   ~\_   _-~            /\
 *                            /  \     \__   \/~                \__
 *                        _--~ _/ | .-~~____--~-/                  ~~==.
 *                       ((->/~   '.|||' -_|    ~~-/ ,              . _||
 *                                  -_     ~\      ~~---l__i__i__i--~~_/
 *                                  _-~-__   ~)  \--______________--~~
 *                                //.-~~~-~_--~- |-------~~~~~~~~
 *                                       //.-~~~--\
 *                       ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 *
 *                               神兽保佑            永无BUG
 *
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2024-05-20 07:15:53
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2024-05-20 07:23:48
 * @FilePath: /github.com/rz1998/invest/investBasic/utilTime.go
 * @Description:
 *
 */

package investBasic

import "time"

func GetTarTime(strDate, strTime string) time.Time {
	time.Local = time.FixedZone("CST", 8*3600)
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", strDate+" "+strTime, time.Local)
	return t
}

// DaysBetweenDates 计算两个日期之间的天数
func DaysBetweenDates(dateBegin, dateEnd string) int64 {
	timeLayout := "2006-01-02"
	// 转换为时间
	begin, _ := time.Parse(timeLayout, dateBegin)
	end, _ := time.Parse(timeLayout, dateEnd)
	// 转换为时间戳
	timeBegin := begin.Unix()
	timeEnd := end.Unix()
	return (timeEnd - timeBegin) / 86400
}
