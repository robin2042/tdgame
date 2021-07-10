package games

import "testing"

func TestSortCardList(t *testing.T) {

	cbCardData := [5]byte{50, 59, 27, 41, 61}
	SortCardList(&cbCardData, 5)

}
