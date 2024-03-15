/*
 * author: I Komang Mardika
 * email: komang.mardika@hotmail.com
 * description: mini project 2 for markasbali ft kominfo golang hacker workshop
 *
 */
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type Book struct {
	Code      string `json:"code"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Publisher string `json:"publisher"`
	Year      int    `json:"year"`
	PageTotal int    `json:"page_total"`
}

var allBooks []Book
var pdfMutex sync.Mutex

/*
 * main functions:
 * show all
 * add new book
 * delete a book
 * edit a book
 * print book(s)
 */

func process(jsonChan chan []byte, errChan chan error, wg *sync.WaitGroup) {
	for jsonData := range jsonChan {
		var b Book
		err := json.Unmarshal(jsonData, &b)
		if err != nil {
			errChan <- err
		}
		filename := fmt.Sprintf("book-%s.json", b.Code)
		if err := saveJsonToFile(jsonData, filename); err != nil {
			errChan <- err
		}
	}
	wg.Done()
}
func convertToJson(book Book) ([]byte, error) {
	encoded, err := json.MarshalIndent(book, "", "    ")
	if err != nil {
		fmt.Printf("")
	}

	return encoded, err
}
func saveJsonToFile(encoded []byte, fileName string) error {
	create, err := os.Create(fmt.Sprintf("books/book-%s.json", fileName))
	if err != nil {
		fmt.Printf("")
	}

	defer func() {
		err := create.Close()
		if err != nil {
			return
		}
	}()

	_, err = create.Write(encoded)
	if err != nil {
		fmt.Printf("")
	}

	return err

}

func showAll(print bool) {

	repopulateBooks()

	if print {
		if len(allBooks) > 0 {
			for i := 0; i < len(allBooks); i++ {
				fmt.Println("No.", i+1)
				fmt.Printf(`
	Code: %s
	Title: %s
	Author: %s
	Publisher: %s
	Publication Year: %d
	Number of Pages: %d
`,
					allBooks[i].Code,
					allBooks[i].Title,
					allBooks[i].Author,
					allBooks[i].Publisher,
					allBooks[i].Year,
					allBooks[i].PageTotal,
				)
			}
		} else {
			fmt.Println("No Book in directory - please add new book first")
		}

	}
}
func repopulateBooks() {
	var wg sync.WaitGroup
	bookCh := make(chan Book)

	files, err := os.ReadDir("books")
	if err != nil {
		fmt.Println("Error reading directory contents:", err)
		return
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			wg.Add(1)

			go func(filename string) {
				defer wg.Done()

				filePath := filepath.Join("books", filename)
				file, err := os.Open(filePath)
				if err != nil {
					fmt.Printf("Error opening file %s: %v\n", filePath, err)
					return
				}
				defer func(file *os.File) {
					err := file.Close()
					if err != nil {

					}
				}(file)

				book, err := convertFileToJson(file)
				bookCh <- book
			}(file.Name())
		}
	}

	go func() {
		wg.Wait()
		close(bookCh)
	}()

	allBooks = nil
	for book := range bookCh {
		allBooks = append(allBooks, book)
	}
}
func addNewBook() {
	var bookId, bookTitle, bookAuthor, bookPublisher string
	var bookYear, bookPageNum int
	var books []Book

	for {
		for {
			bookId = inputText("Enter new Book Code: ")
			bookExists := isBookExistsByCode(bookId)

			if !bookExists {
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

		if keyPress("Press y to add another book ") {
			continue
		} else {
			break
		}
	} //end for

	clearScreen()
	fmt.Println("Saving all of your book entries")

	var wg sync.WaitGroup
	wg.Add(len(books))

	errChannel := make(chan error, len(books))
	jsonChannel := make(chan []byte, len(books))

	const numWorkers = 2
	for i := 0; i < numWorkers; i++ {
		go func() {
			for _, book := range books {
				currentBook := book
				jsonData, err := convertToJson(currentBook)
				if err != nil {
					errChannel <- err
					return
				}

				err = saveJsonToFile(jsonData, currentBook.Code)
				if err != nil {
					errChannel <- err
					return
				}
			}

		}()
	}

	go process(jsonChannel, errChannel, &wg)

	fmt.Println("Book has been added successfully")
	confirmation()
}
func deleteBook() {
	bookCode := inputText("Enter book code you want to delete: ")
	exists := isBookExistsByCode(bookCode)
	if exists {
		err := os.Remove(fmt.Sprintf("books/book-%s.json", bookCode))
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	} else {
		fmt.Printf("book with code %s is not found \n", bookCode)
		return
	}
	clearScreen()
	fmt.Println("New Books List after Deletion")
	showAll(true)

}
func editBook() {
	bookCode := inputText("Enter book code you want to edit: ")
	b, err := findBookJsonFileByCode(bookCode)
	if err == nil {
		fmt.Println("Book Code: ", b.Code)
		b.Title = inputText("Enter new Title: ")
		b.Author = inputText("Enter new Author: ")
		b.Publisher = inputText("Enter new Publisher: ")
		b.Year = inputInt("Enter new Year: ")
		b.PageTotal = inputInt("Enter new Total Page: ")

		toJson, err := convertToJson(b)
		if err != nil {
			return
		}

		err = saveJsonToFile(toJson, b.Code)
		if err != nil {
			return
		}

	} else {
		fmt.Printf("book with code %s not found \n", bookCode)
		return
	}
	clearScreen()
	fmt.Println("New Books List after Alteration")
	showAll(true)
}
func printBook() {
	if keyPress("Press y to print all books ") {
		printAllBook()
	} else {
		bookId := inputText("Enter new Book Code: ")
		book, err := findBookJsonFileByCode(bookId)
		if err != nil {
			fmt.Println("error book not found")
			return
		}

		printSingleBook(book)

	}
}
func printSingleBook(book Book) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "", 12)

	writeHeadingToPdf(pdf)
	writeBookToPDF(pdf, book)

	err := pdf.OutputFileAndClose(fmt.Sprintf("pdf/%s.pdf", book.Code))
	if err != nil {
		fmt.Println("Error saving PDF file:", err)
		return
	}
}
func printAllBook() {
	repopulateBooks()

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "", 10)
	writeHeadingToPdf(pdf)

	bookChan := make(chan Book)
	done := make(chan bool)

	go func() {
		for _, book := range allBooks {
			bookChan <- book
		}
		close(bookChan)
	}()

	var wg sync.WaitGroup

	const numWorkers = 2
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for book := range bookChan {
				writeBookToPDF(pdf, book)
			}
		}()
	}

	// Close 'done' channel when all goroutines are finished
	go func() {
		wg.Wait()
		close(done)
	}()

	<-done

	err := pdf.OutputFileAndClose("pdf/all-books.pdf")
	if err != nil {
		fmt.Println("Error saving PDF:", err)
		return
	}
}
func writeHeadingToPdf(pdf *gofpdf.Fpdf) {
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(20, 10, "Code", "1", 0, "", true, 0, "")
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(40, 10, "Title", "1", 0, "", true, 0, "")
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(40, 10, "Author", "1", 0, "", true, 0, "")
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(40, 10, "Publisher", "1", 0, "", true, 0, "")
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(20, 10, "Year", "1", 0, "R", true, 0, "")
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(20, 10, "Page", "1", 0, "R", true, 0, "")
	pdf.Ln(-1)
}
func writeBookToPDF(pdf *gofpdf.Fpdf, book Book) {
	pdfMutex.Lock()
	defer pdfMutex.Unlock()
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(20, 10, book.Code, "1", 0, "", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(40, 10, book.Title, "1", 0, "", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(40, 10, book.Author, "1", 0, "", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(40, 10, book.Publisher, "1", 0, "", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(20, 10, fmt.Sprintf("%d", book.Year), "1", 0, "R", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(20, 10, fmt.Sprintf("%d", book.PageTotal), "1", 0, "R", false, 0, "")
	pdf.Ln(-1)
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
	5. print book(s)
	6. Quit
`)
	menuChosen = inputText("Please choose your option(1/2/3/4/5): ")

	switch menuChosen {
	case "1":
		showAll(true)
	case "2":
		addNewBook()
	case "3":
		deleteBook()
	case "4":
		editBook()
	case "5":
		printBook()
	case "6":
		fmt.Println("Thanks for using Simple Library App")
		os.Exit(0)
	default:
		fmt.Print("You did not choose the right option, ")
	}

	confirmation()
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
	if keyPress("Press y to return to the menu ") {
		menu()
	}
}
func keyPress(msg string) bool {
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
func convertFileToJson(file *os.File) (Book, error) {
	bytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error:", err)
	}

	var data Book

	err = json.Unmarshal(bytes, &data)
	return data, err
}
func findBookJsonFileByCode(code string) (Book, error) {
	file, err := os.Open(fmt.Sprintf("books/book-%s.json", code))
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error:", err)
		}
	}(file)

	return convertFileToJson(file)
}
func isBookExistsByCode(code string) bool {
	_, err := os.Stat(fmt.Sprintf("books/book-%s.json", code))

	// Check if the file exists
	if os.IsNotExist(err) {
		return false
	}
	return true

}
