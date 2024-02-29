/*
 * author: I Komang Mardika
 * email: komang.mardika@hotmail.com
 * description: mini project 1 for markasbali ft kominfo golang hacker workshop
 *
 */
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Book struct {
	Code      string
	Title     string
	Author    string
	Publisher string
	Year      int
	PageTotal int
}

/*
 * main functions:
 * show all
 * add new book
 * delete a book
 * edit a book
 */
func showAll() {
	for i := 0; i < len(books); i++ {
		fmt.Println("No.", i+1)
		fmt.Printf(`
	Code: %s
	Title: %s
	Author: %s
	Publisher: %s
	Publication Year: %d
	Number of Pages: %d
`,
			books[i].Code,
			books[i].Title,
			books[i].Author,
			books[i].Publisher,
			books[i].Year,
			books[i].PageTotal,
		)
	}
}
func addNewBook() {
	var bookId, bookTitle, bookAuthor, bookPublisher string
	var bookYear, bookPageNum int

	for {
		bookId = inputText("Enter new Book Code: ")
		res := findBookIndexByCode(bookId)
		fmt.Println(res)

		if res == -1 {
			break
		} else {
			fmt.Println("Book code is exists, try different book code")
			continue
		}
	}

	bookTitle = inputText("Enter new title: ")
	bookAuthor = inputText("Enter the author: ")
	bookPublisher = inputText("Enter the publisher: ")
	bookYear = inputInt("Enter Publication Year: ")
	bookPageNum = inputInt("Enter Book Total Page: ")

	books = append(books, Book{
		Code:      bookId,
		Title:     bookTitle,
		Author:    bookAuthor,
		Publisher: bookPublisher,
		Year:      bookYear,
		PageTotal: bookPageNum,
	})
	fmt.Println("Book has been added successfully")
	confirmation()
}
func deleteBook() {
	bookCode := inputText("Enter book code you want to delete: ")
	index := findBookIndexByCode(bookCode)
	if index > -1 {
		books = append(books[:index], books[index+1:]...)
	} else {
		fmt.Printf("book with code %s is not found \n", bookCode)
		return
	}
	clearScreen()
	fmt.Println("New Books List after Deletion")
	showAll()

}
func editBook() {
	bookCode := inputText("Enter book code you want to edit: ")
	index := findBookIndexByCode(bookCode)
	if index > -1 {
		fmt.Println("Book Code: ", books[index].Code)
		books[index].Title = inputText("Enter new Title: ")
		books[index].Author = inputText("Enter new Author: ")
		books[index].Publisher = inputText("Enter new Publisher: ")
		books[index].Year = inputInt("Enter new Year: ")
		books[index].PageTotal = inputInt("Enter new Total Page: ")

	} else {
		fmt.Printf("book with code %s not found \n", bookCode)
		return
	}
	clearScreen()
	fmt.Println("New Books List after Alteration")
	showAll()
}

/* menu */
func menu() {
	var menuChosen string

	clearScreen()

	fmt.Printf(`
Simple Library App 
	1. show books
	2. add a new book
	3. delete a book
	4. edit a book
	5. Quit
`)
	menuChosen = inputText("Please choose your option(1/2/3/4): ")

	switch menuChosen {
	case "1":
		showAll()
	case "2":
		addNewBook()
	case "3":
		deleteBook()
	case "4":
		editBook()
	case "5":
		fmt.Println("Thanks for using Simple Library App")
		os.Exit(0)
	default:
		fmt.Print("You did not choose the right option, ")
	}

	confirmation()
}

/* db seeder */
var books = []Book{
	{
		Code:      "B-0001",
		Title:     "Harry Potter and Sorcerer's Stone",
		Author:    "JK Rowling",
		Publisher: "Gramedia",
		Year:      2000,
		PageTotal: 10,
	},
	{
		Code:      "B-0002",
		Title:     "Harry Potter and Chamber of Secret",
		Author:    "JK Rowling",
		Publisher: "Gramedia",
		Year:      2001,
		PageTotal: 10,
	},
	{
		Code:      "B-0003",
		Title:     "Harry Potter and Prisoner of Azkaban",
		Author:    "JK Rowling",
		Publisher: "Gramedia",
		Year:      2002,
		PageTotal: 10,
	},
}

/* runner */
func main() {
	menu()
}

/*
 * utilities:
 * clear console screen
 * confirmation for repeat
 * detect keypress of confirmation
 */
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
func confirmation() {
	if keyPress() {
		menu()
	}
}
func keyPress() bool {
	fmt.Print("Press y to return to the menu ")
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

/*
 * helpers:
 * input text
 * input int
 * find book index by code
 * generate book code
 */
func inputText(text string) string {
	var typedText string
	fmt.Print(text)

	readerTitle := bufio.NewReader(os.Stdin)
	typedText, _ = readerTitle.ReadString('\n')
	return strings.TrimSpace(typedText)
}
func inputInt(text string) int {
	var typedText string
	var number int
	for {
		typedText = inputText(text)
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
func findBookIndexByCode(code string) int {
	index := -1
	for i := 0; i < len(books); i++ {
		if code == books[i].Code {
			index = i
			break
		}
	}
	return index
}
func generateBookCode(lastCode string) string {
	parts := strings.Split(lastCode, "-")
	number, _ := strconv.Atoi(parts[1])
	number = number + 1

	if number < 10 {
		return "B-000" + strconv.Itoa(number)
	} else if number > 9 && number <= 99 {
		return "B-00" + strconv.Itoa(number)
	} else if number > 99 && number <= 999 {
		return "B-0" + strconv.Itoa(number)
	} else {
		return "B-" + strconv.Itoa(number)
	}
}
