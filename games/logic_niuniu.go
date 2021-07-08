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
	LOGIC_MASK_COLOR = 0xF0 //花色掩码
	LOGIC_MASK_VALUE = 0x0F //数值掩码
)

var m_cbTableCardArray [5][5]byte //牌

var m_cbCardListData = [...]byte{
	0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, //方块 A - K
	0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, //梅花 A - K
	0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D, //红桃 A - K
	0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D, //黑桃 A - K
	0x4E, 0x4F,
}

//generateRandomNumber	 生成随机数
func generateRandomNumber(start int, end int, count int) []int {
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
