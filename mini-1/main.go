package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Book struct {
	Id     string
	Title  string
	Status string
}

func transactionBook(action string) {
	var bookId string
	fmt.Print("Masukkan ID buku yang ingin ")

	if action == "B" {
		fmt.Print("dipinjam : ")
	} else {
		fmt.Print("dikembalikan : ")
	}

	reader := bufio.NewReader(os.Stdin)
	bookId, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	} else {
		bookId = strings.TrimSpace(bookId)

		flag := false
		for _, book := range books {
			if book.Id == bookId {
				fmt.Println("buku ditemukan berikut data nya")
				fmt.Println("Book ID: ", book.Id)
				fmt.Println("Book Title: ", book.Title)
				fmt.Println("status buku sudah diubah")
				if action == "B" {
					book.Status = "B"
				} else {
					book.Status = "A"
				}
				flag = true
				break
			}
		}

		if !flag {
			confirmation()
		} else {
			println("buku tidak ditemukan")
			confirmation()
		}
	}
}
func showAll(status string) {
	var s string
	for i := 0; i < len(books); i++ {
		if status == books[i].Status {
			fmt.Printf(`
	id: %s
	title: %s
`, books[i].Id, books[i].Title)
		} else {
			if status == "all" {
				if books[i].Status == "A" {
					s = "Tersedia"
				} else {
					s = "Dipinjam"
				}
				fmt.Printf(`
	id: %s
	title: %s
	status: %s
`, books[i].Id, books[i].Title, s)
			}
		}

	}
}
func addNewBook() {
	var bookId string

	fmt.Print("Masukkan Title buku baru : ")

	readerTitle := bufio.NewReader(os.Stdin)
	bookTitle, _ := readerTitle.ReadString('\n')

	index := len(books) - 1
	parts := strings.Split(books[index].Id, "-")
	number, _ := strconv.Atoi(parts[1])
	number = number + 1

	if number < 10 {
		bookId = "B-000" + strconv.Itoa(number)
	} else if number > 9 && number <= 99 {
		bookId = "B-00" + strconv.Itoa(number)
	} else if number > 99 && number <= 999 {
		bookId = "B-0" + strconv.Itoa(number)
	} else {
		bookId = "B-" + strconv.Itoa(number)
	}

	books = append(books, Book{Id: bookId, Title: bookTitle, Status: "A"})
	fmt.Println("Buku berhasil ditambahkan")
	confirmation()
}
func menu() {
	var menuChosen string

	clearScreen()

	fmt.Printf(`
Simple Perpustakaan App 
	1. lihat buku tersedia
	2. lihat buku dipinjam
	3. lihat seluruh buku
	4. menambah buku
	5. pinjam buku
	6. kembalikan buku
	7. Keluar dari aplikasi
`)
	fmt.Print("Pilihan anda: ")
	reader := bufio.NewReader(os.Stdin)
	menuChosen, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	} else {

		menuChosen = strings.TrimSpace(menuChosen)

		switch menuChosen {
		case "1":
			showAll("A")
		case "2":
			showAll("B")
		case "3":
			showAll("all")
		case "4":
			addNewBook()
		case "5":
			transactionBook("B")
		case "6":
			transactionBook("R")
		case "7":
			fmt.Println("Terima kasih sudah menggunakan aplikasi ini")
			return
		default:
			fmt.Print("Anda tidak memilih menu yang tepat, ")
		}

		confirmation()
	}
}
func confirmation() {
	if keyPress() {
		menu()
	} else {
		fmt.Println("Terima kasih sudah menggunakan aplikasi ini")
		return
	}
}
func keyPress() bool {
	fmt.Println("tekan y untuk mengulang")
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

var books = []Book{
	{
		Id:     "B-00001",
		Title:  "Harry Potter and Sorcerer's Stone",
		Status: "A",
	},
	{
		Id:     "B-00002",
		Title:  "Harry Potter and Chamber of Secret",
		Status: "A",
	},
	{
		Id:     "B-00003",
		Title:  "Harry Potter and Prisoner of Azkaban",
		Status: "A",
	},
}

func main() {
	// seed books
	menu()
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
