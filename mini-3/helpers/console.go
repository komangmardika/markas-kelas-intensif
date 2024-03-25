package helpers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

func KeyPress(msg string) bool {
	fmt.Print(msg)
	reader := bufio.NewReader(os.Stdin)
	for {
		key, _, _ := reader.ReadRune()
		if key == 'y' {
			return true
		} else {
			return false
		}
	}
}

func InputText(text string) string {
	var typedText string
	fmt.Print(text)

	readerTitle := bufio.NewReader(os.Stdin)
	typedText, _ = readerTitle.ReadString('\n')
	return strings.TrimSpace(typedText)
}
func InputInt(text string) int {
	var typedText string
	var number int
	for {
		typedText = InputText(text)
		num, err := strconv.Atoi(strings.TrimSpace(typedText))
		if err != nil {
			fmt.Println("you are not entering number")
			continue
		} else {
			number = num
			break
		}
	}

	return number
}

func InputUint(text string) uint {
	var typedText string
	var number uint
	for {
		typedText = InputText(text)
		num, err := strconv.ParseUint(strings.TrimSpace(typedText), 10, 64)
		if err != nil {
			fmt.Println("you are not entering number")
			continue
		} else {
			number = uint(num)
			break
		}
	}

	return number
}
