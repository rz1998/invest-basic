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
 * @Date: 2024-05-16 01:39:05
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2024-06-01 01:13:17
 * @FilePath: /investBasic/test/logic_test.go
 * @Description:
 *
 */

package test

import (
	"testing"
	"time"

	investBasic "github.com/rz1998/invest-basic"
)

// 测试主力合约计算
func TestQryMainContract(t *testing.T) {
	product := "IC"
	date, _ := time.Parse("2006-01-02", "2024-05-15")
	mainContract := investBasic.QryMainContract(product, date)
	if mainContract != "IC2406.CFFEX" {
		t.Errorf("%s %s's main contract expected to be IC2406.CFFEX, but got %s", product, date.Format("2006-01-02"), mainContract)
	}
	date, _ = time.Parse("2006-01-02", "2024-05-31")
	mainContract = investBasic.QryMainContract(product, date)
	if mainContract != "IC2406.CFFEX" {
		t.Errorf("%s %s's main contract expected to be IC2406.CFFEX, but got %s", product, date.Format("2006-01-02"), mainContract)
	}
}

// 测试中位数计算
func TestCalMedian(t *testing.T) {
	array1 := []float64{1, 2, 3, 4, 5}
	median1 := investBasic.CalMedian(array1)
	array2 := []float64{1, 2, 3, 4}
	median2 := investBasic.CalMedian(array2)
	if median1 != 3 || median2 != 2.5 {
		t.Errorf("Median calculation failed")
	}
}

func TestTime(t *testing.T) {
	strMethod := "TestTime"
	// GetTarTime
	tarTime := investBasic.GetTarTime("2000-01-01", "09:00:00")
	if tarTime.Format("2006-01-02 15:04:05") != "2000-01-01 09:00:00" {
		t.Errorf("%s func error %s", strMethod, "GetTarTime")
	}
}
