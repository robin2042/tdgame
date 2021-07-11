package games

import (
	"fmt"
	"math/rand"
	"time"
)

//数值掩码

const (
	GAME_PLAYER      = 5
	MAX_COUNT        = 5
	MAX_MULTIPLE     = 10   //最高倍数
	LOGIC_MASK_COLOR = 0xF0 //花色掩码
	LOGIC_MASK_VALUE = 0x0F //数值掩码
	//扑克类型
	OX_VALUE0 = 0 //混合牌型  无牛

	OX_SMALL_WANG = 11 //小王牛
	OX_BIG_WANG   = 12 //大王牛
	OX_DOUBLECOW  = 13 //牛牛

	OX_FOUR_SAME = 15 //炸弹牌型
	OX_FOURKING  = 14 //四花牛牌型
	OX_FIVEKING  = 16 //五花牛牌型
	OX_FIVESMALL = 17 //五小
)

var m_cbCardListData = [...]byte{
	0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, //方块 A - K
	0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, //梅花 A - K
	0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D, //红桃 A - K
	0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D, //黑桃 A - K
	// 0x4E, 0x4F,
}

//逻辑数值
func GetCardLogicValue(cbCardData byte) int {
	//扑克属性
	// bCardColor := GetCardColor(cbCardData)
	bCardValue := GetCardValue(cbCardData)

	//转换数值
	return If(int(bCardValue) > 10, 10, int(bCardValue)).(int)

}

//生成count个[start,end)结束的不重复的随机数

func GenerateRandomNumber(start int, end int, count int) []int {

	//范围检查

	if end < start || (end-start) < count {

		return nil

	}

	//存放结果的slice

	nums := make([]int, 0)

	//随机数生成器，加入时间戳保证每次生成的随机数不一样

	r := rand.New(rand.NewSource(time.Now().UnixNano() + time.Now().UnixNano()))

	for len(nums) < count {

		//生成随机数

		num := r.Intn((end - start)) + start

		//查重

		exist := false

		for _, v := range nums {

			if v == num {

				exist = true

				break

			}

		}

		if !exist {

			nums = append(nums, num)

		}

	}

	return nums

}

//获取数值
func GetCardValue(cbCardData byte) byte { return cbCardData & LOGIC_MASK_VALUE }

//获取花色
func GetCardColor(cbCardData byte) byte { return cbCardData & LOGIC_MASK_COLOR }

func GetUnicode(cbCardData [5]byte, cbCardCount byte) {
	var str string
	for i := 0; i < len(cbCardData); i++ {

		if cbCardData[0] == 0x4E {
			str += "小王"
			continue
		} else if cbCardData[i] == 0x4F {
			str += " 大王"
			continue
		}

		switch GetCardColor(cbCardData[i]) {
		case 0x00:
			str += " 方块"
			break
		case 0x10:
			str += " 梅花"
			break
		case 0x20:
			str += " 红桃"
			break
		case 0x30:
			str += " 黑桃"
			break
		}
		switch GetCardValue(cbCardData[i]) {
		case 1:
			str += "A"
			break
		case 2:
			str += "2"
			break
		case 3:
			str += "3"
			break
		case 4:
			str += "4"
			break
		case 5:
			str += "5"
			break
		case 6:
			str += "6"
			break
		case 7:
			str += "7"
			break
		case 8:
			str += "8"
			break
		case 9:
			str += "9"
			break
		case 10:
			str += "10"
			break
		case 11:
			str += "J"
			break
		case 12:
			str += "Q"
			break
		case 13:
			str += "K"
			break
		}
	}
	fmt.Println(str)
}

// //获取类型
func GetCardType(cbCardData [MAX_COUNT]byte, cbCardCount int) int {
	var bKingCount, bTenCount int
	var bSmallWang, bBigWang bool
	for i := 0; i < cbCardCount; i++ {
		//大小王可以变花牌
		if cbCardData[i] == 0x4E {
			bSmallWang = true
			bKingCount++
			continue
		}
		if cbCardData[i] == 0x4F {
			bBigWang = true
			bKingCount++
			continue
		}
		if GetCardValue(cbCardData[i]) > 10 {
			bKingCount++
		} else if GetCardValue(cbCardData[i]) == 10 {
			bTenCount++
		}

	}
	if bKingCount == MAX_COUNT {
		return OX_FIVEKING
	}

	bFirstTemp := cbCardData

	SortCardList(&bFirstTemp, cbCardCount)
	//fmt.Println(bSmallWang, bBigWang)

	var TempSum int
	for i := 0; i < cbCardCount; i++ {

		if cbCardData[i] == 0x4E {
			TempSum++
			continue
		}
		if cbCardData[i] == 0x4F {
			TempSum++
			continue
		}
		TempSum += GetCardLogicValue(cbCardData[i])
	}
	if TempSum <= 10 {
		return OX_FIVESMALL
	}

	// 	//有大小王的情况
	if bSmallWang || bBigWang {
		//大小王都有，那就只要判断有两张一样的牌
		if bSmallWang && bBigWang {
			if GetCardValue(bFirstTemp[2]) == GetCardValue(bFirstTemp[3]) || GetCardValue(bFirstTemp[3]) == GetCardValue(bFirstTemp[4]) {
				return OX_FOUR_SAME
			}
		} else ////只有一个大小王，那就判断三张一样的牌
		{
			if GetCardValue(bFirstTemp[1]) == GetCardValue(bFirstTemp[3]) || GetCardValue(bFirstTemp[2]) == GetCardValue(bFirstTemp[4]) {
				return OX_FOUR_SAME
			}
		}
	} else //没有大小王的情况，判断四张
	{
		if GetCardValue(bFirstTemp[0]) == GetCardValue(bFirstTemp[3]) || GetCardValue(bFirstTemp[1]) == GetCardValue(bFirstTemp[4]) {
			return OX_FOUR_SAME
		}
	}

	if bKingCount == MAX_COUNT-1 && bTenCount == 1 {
		return OX_FOURKING
	}

	if bSmallWang || bBigWang {
		if bSmallWang && bBigWang {
			return OX_DOUBLECOW
		}
		var cbMaxTwoCardValue, cbTemp int
		var cbTempValue [3]int
		for i := 0; i < cbCardCount; i++ {

			if cbCardData[i] == 0x4E || cbCardData[i] == 0x4F {
				continue
			}
			for j := 0; j < cbCardCount; j++ {

				if cbCardData[j] == 0x4E || cbCardData[j] == 0x4F {
					continue
				}
				cbTempValue[0] = GetCardLogicValue(cbCardData[i])
				cbTempValue[1] = GetCardLogicValue(cbCardData[j])

				cbTemp = If(cbTempValue[0]+cbTempValue[1] > 10, cbTempValue[0]+cbTempValue[1]-10, cbTempValue[0]+cbTempValue[1]).(int)

				if cbMaxTwoCardValue < cbTemp {
					cbMaxTwoCardValue = cbTemp
				}
			}
		}
		for i := 0; i < cbCardCount; i++ {

			if cbCardData[i] == 0x4E || cbCardData[i] == 0x4F {
				continue
			}
			for j := 0; j < cbCardCount; j++ {

				if cbCardData[j] == 0x4E || cbCardData[j] == 0x4F {
					continue
				}

				for k := 0; k < cbCardCount; k++ {

					if cbCardData[k] == cbCardData[i] || cbCardData[k] == cbCardData[j] || cbCardData[k] == 0x4E || cbCardData[k] == 0x4F {
						continue
					}
					cbTempValue[0] = GetCardLogicValue(cbCardData[i])
					cbTempValue[1] = GetCardLogicValue(cbCardData[j])
					cbTempValue[2] = GetCardLogicValue(cbCardData[k])
					if (cbTempValue[0]+cbTempValue[1]+cbTempValue[2])%10 == 0 {
						return If(bSmallWang == true, OX_DOUBLECOW, OX_DOUBLECOW).(int)

					}
				}
			}
		}
		if cbMaxTwoCardValue == 10 {
			return If(bSmallWang == true, OX_DOUBLECOW, OX_DOUBLECOW).(int)
		} else {
			return cbMaxTwoCardValue
		}
	}

	var bTemp [MAX_COUNT]int
	var bSum int
	for i := 0; i < cbCardCount; i++ {

		bTemp[i] = GetCardLogicValue(cbCardData[i])
		bSum += bTemp[i]
	}
	var cbTempValue int
	for i := 0; i < cbCardCount; i++ {

		for j := 0; j < cbCardCount; j++ {

			if (bSum-bTemp[i]-bTemp[j])%10 == 0 {

				cbTempValue = If((bTemp[i]+bTemp[j]) > 10, bTemp[i]+bTemp[j]-10, bTemp[i]+bTemp[j]).(int)
				if cbTempValue == 10 {
					return OX_DOUBLECOW
				} else {
					return cbTempValue
				}
			}
		}
	}
	return OX_VALUE0

}

// //对比扑克
func CompareCard(cbFirstData [MAX_COUNT]byte, cbNextData [MAX_COUNT]byte, cbCardCount int) bool {
	// 	//获取点数
	cbNextType := GetCardType(cbNextData, cbCardCount)
	cbFirstType := GetCardType(cbFirstData, cbCardCount)
	//点数判断
	if cbFirstType != cbNextType {
		return (cbFirstType > cbNextType)
	}

	//排序大小
	// var bFirstTemp [MAX_COUNT]byte
	// var bNextTemp [MAX_COUNT]byte
	bFirstTemp := cbFirstData
	bNextTemp := cbNextData
	SortCardList(&bFirstTemp, cbCardCount)
	SortCardList(&bNextTemp, cbCardCount)

	// 	//比较数值
	cbNextMaxValue := GetCardValue(bNextTemp[0])
	cbFirstMaxValue := GetCardValue(bFirstTemp[0])
	var cbPosFirst, cbPosNext byte
	if cbNextMaxValue != cbFirstMaxValue {
		return cbFirstMaxValue > cbNextMaxValue
	}
	// 	//比较颜色
	return GetCardColor(bFirstTemp[cbPosFirst]) > GetCardColor(bNextTemp[cbPosNext])

}

func SortCardList(cbCardData *[5]byte, cbCardCount int) {
	//转换数值
	var cbLogicValue [MAX_COUNT]byte
	for i := 0; i < cbCardCount; i++ {
		cbLogicValue[i] = GetCardValue(cbCardData[i])
	}

	//排序操作
	var bSorted bool
	var cbTempData byte
	var bLast = cbCardCount - 1
	for {
		bSorted = true
		for i := 0; i < bLast; i++ {
			if (cbLogicValue[i] < cbLogicValue[i+1]) ||
				((cbLogicValue[i] == cbLogicValue[i+1]) && (cbCardData[i] < cbCardData[i+1])) {
				//交换位置
				cbTempData = cbCardData[i]
				cbCardData[i] = cbCardData[i+1]
				cbCardData[i+1] = cbTempData
				cbTempData = cbLogicValue[i]
				cbLogicValue[i] = cbLogicValue[i+1]
				cbLogicValue[i+1] = cbTempData
				bSorted = false
			}
		}
		bLast--
		if bSorted {
			break

		}

	}
}

// //排列扑克
// func SortCardList(cbCardData *[5]byte, cbCardCount int) {
// 	//转换数值
// 	var cbLogicValue [MAX_COUNT]byte
// 	for i := 0; i < cbCardCount; i++ {
// 		cbLogicValue[i] = GetCardValue(cbCardData[i])
// 	}

// 	//排序操作
// 	var bSorted bool
// 	var cbTempData byte
// 	var bLast = cbCardCount - 1

// 	for {
// 		bSorted = true
// 		for i := 0; i < bLast; i++ {

// 			if (cbLogicValue[i] < cbLogicValue[i+1]) ||
// 				((cbLogicValue[i] == cbLogicValue[i+1]) && (cbCardData[i] < cbCardData[i+1])) {
// 				//交换位置
// 				cbTempData = cbCardData[i]
// 				cbCardData[i] = cbCardData[i+1]
// 				cbCardData[i+1] = cbTempData
// 				cbTempData = cbLogicValue[i]
// 				cbLogicValue[i] = cbLogicValue[i+1]
// 				cbLogicValue[i+1] = cbTempData
// 				bSorted = false
// 			}
// 		}
// 		bLast--
// 		if !bSorted {
// 			return
// 		}
// 	}

// }
func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

//获取倍数
func GetTimes(cbCardData [MAX_COUNT]byte, cbCardCount byte, lMultiple int) int {
	if cbCardCount != MAX_COUNT {
		return 0
	}

	bTimes := GetCardType(cbCardData, MAX_COUNT)
	if lMultiple == 10 {
		if bTimes < 2 {
			return 1
		}
		if bTimes >= 2 && bTimes <= 9 {
			return bTimes
		} else if bTimes >= OX_DOUBLECOW {
			return 10
		}
	} else if lMultiple == 4 {
		if bTimes < 7 {
			return 1
		}
		if bTimes >= 7 && bTimes <= 8 {
			return 2
		} else if bTimes == 9 {
			return 3
		} else if bTimes >= OX_DOUBLECOW {
			return 4
		}
	}

	return 0
}

//获取花色
func GetCardColorEmoj(cbCardData byte) string {
	var str string
	//var card string
	value := int(GetCardValue(cbCardData))
	card := fmt.Sprintf("%d", value)

	if value == 11 {
		card = "J"
	} else if value == 12 {
		card = "Q"
	} else if value == 13 {
		card = "K"
	}

	switch GetCardColor(cbCardData) {
	case 0x00:
		str = "♦️"

	case 0x10:
		str = "♣"

	case 0x20:
		str = "♥️"

	case 0x30:
		str = "♠️"

	}
	str = fmt.Sprintf("%s%s", card, str)
	return str
}

//获取牌
func GetCardTimesEmoj(cbCardData [MAX_COUNT]byte) string {
	var str string
	value := GetTimes(cbCardData, MAX_COUNT, MAX_MULTIPLE)
	if value == 0 {
		str = "(无牛)"
	} else if value == 1 {
		str = "(牛一)"
	} else if value == 2 {
		str = "(牛二)"
	} else if value == 3 {
		str = "(牛三)"
	} else if value == 4 {
		str = "(牛四)"
	} else if value == 5 {
		str = "(牛五)"
	} else if value == 6 {
		str = "(牛六)"
	} else if value == 7 {
		str = "(牛七)"
	} else if value == 8 {
		str = "(牛八)"
	} else if value == 9 {
		str = "(牛九)"
	} else if value == 10 {
		str = "(牛牛)"
	}

	return str
}

//获取花色
func GetCardValueEmoj(cbCardData [MAX_COUNT]byte) string {
	var str string
	for i := 0; i < MAX_COUNT; i++ {
		str += GetCardColorEmoj(cbCardData[i])
	}
	return str
}
