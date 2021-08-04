package games

const (
	LOGIC_MASK_COLOR = 0xF0 //花色掩码
	LOGIC_MASK_VALUE = 0x0F //数值掩码
)

func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

//获取数值
func GetCardValue(cbCardData byte) byte { return cbCardData & LOGIC_MASK_VALUE }

//获取花色
func GetCardColor(cbCardData byte) byte { return cbCardData & LOGIC_MASK_COLOR }
