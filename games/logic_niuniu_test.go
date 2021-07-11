package games

import (
	"fmt"
	"testing"
)

func TestSortCardList(t *testing.T) {

	cbCardData := [5]byte{17, 23, 5, 28, 26}

	// cbCardData := [5]byte{50, 59, 27, 41, 61}
	bank := [5]byte{39, 25, 29, 53, 58}

	fmt.Println(CompareCard(cbCardData, bank, 5))

}

//61 59 27 41 50

func TestGetCardColorEmoj(t *testing.T) {
	fmt.Println(GetCardColorEmoj(26))
}

func TestGetCardValueEmoj(t *testing.T) {
	cbCardData := [5]byte{17, 23, 5, 28, 26}
	str := GetCardValueEmoj(cbCardData)
	fmt.Println(str)

}

func TestGetCardTimesEmoj(t *testing.T) {
	// cbCardData := [5]byte{17, 23, 5, 28, 26}
	bank := [5]byte{56, 2, 53, 12, 20}

	// str := GetCardValueEmoj(cbCardData)

	fmt.Println(GetCardTimesEmoj(bank))

}
