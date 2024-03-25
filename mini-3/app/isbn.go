package app

import (
	"math/rand"
	"strconv"
)

func GenerateISBN() string {
	var digits []int
	for i := 0; i < 12; i++ {
		digits = append(digits, rand.Intn(10))
	}

	sum := 0
	for i, digit := range digits {
		if i%2 == 0 {
			sum += digit
		} else {
			sum += digit * 3
		}
	}
	checkDigit := (10 - (sum % 10)) % 10

	isbn := ""
	for _, digit := range digits {
		isbn += strconv.Itoa(digit)
	}
	isbn += strconv.Itoa(checkDigit)

	return isbn
}
