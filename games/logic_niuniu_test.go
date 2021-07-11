package games

import (
	"testing"
)

func TestSortCardList(t *testing.T) {

	cbCardData := [5]byte{50, 59, 27, 41, 61}
	bank := [5]byte{39, 25, 29, 53, 58}

	CompareCard(cbCardData, bank, 5)

}

//61 59 27 41 50
