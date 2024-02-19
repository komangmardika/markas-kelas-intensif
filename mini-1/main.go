/*
 * author: I Komang Mardika
 * email: komang.mardika@hotmail.com
 * app desc: mini project 1 for Markas Bali ft Kominfo Kelas Intensif 2024
 * first created on Monday, 19 feb 2024
 */
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/* define default books */
var books = []map[string]string{
	{
		"id":     "B-0001",
		"title":  "Harry Potter And Sorcerer's Stone",
		"status": "A",
	},
	{
		"id":     "B-0002",
		"title":  "Harry Potter And Chamber of Secret",
		"status": "B",
	},
	{
		"id":     "B-0003",
		"title":  "Harry Potter And Prisoner of Azkaban",
		"status": "A",
	},
	{
		"id":     "B-0004",
		"title":  "Harry Potter And Goblet of Fire",
		"status": "A",
	},
}

func main() {
	/* call showMenu function for first time */
	showMenu()
}
func showMenu() {
	var menuChosen string

	clearConsole()

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
			showAllBooks("A")
		case "2":
			showAllBooks("B")
		case "3":
			showAllBooks("all")
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

		confirm()
	}

}

func clearConsole() {
	fmt.Print("\033[H\033[2J")
}

func holdKey() bool {
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

func showAllBooks(param string) {
	var status string
	for i := 0; i < len(books); i++ {
		if param == books[i]["status"] {
			fmt.Printf(`
	id: %s
	title: %s
`, books[i]["id"], books[i]["title"])
		} else {
			if param == "all" {
				if books[i]["status"] == "A" {
					status = "Tersedia"
				} else {
					status = "Dipinjam"
				}
				fmt.Printf(`
	id: %s
	title: %s
	status: %s
`, books[i]["id"], books[i]["title"], status)
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
	parts := strings.Split(books[index]["id"], "-")
	number, _ := strconv.Atoi(parts[1])
	number = number + 1
	if number < 10 {
		bookId = "B-000" + strconv.Itoa(number)
	} else if number > 9 && number < 99 {
		bookId = "B-00" + strconv.Itoa(number)
	} else if number > 99 && number < 999 {
		bookId = "B-0" + strconv.Itoa(number)
	} else {
		bookId = "B-" + strconv.Itoa(number)
	}

	books = append(books, map[string]string{"id": bookId, "title": bookTitle, "status": "A"})
	fmt.Println("Buku berhasil ditambahkan")
	confirm()
}

func transactionBook(param string) {
	var bookId string
	fmt.Print("Masukkan ID buku yang ingin ")
	if param == "B" {
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
			if book["id"] == bookId {
				fmt.Println("buku ditemukan berikut data nya")
				fmt.Println("Book ID: ", book["id"])
				fmt.Println("Book Title: ", book["title"])
				fmt.Println("status buku sudah diubah")
				if param == "B" {
					book["status"] = "B"
				} else {
					book["status"] = "A"
				}
				flag = true
				break
			}
		}

		if !flag {
			confirm()
		} else {
			println("buku tidak ditemukan")
			confirm()
		}
	}
}

func confirm() {
	if holdKey() {
		showMenu()
	} else {
		fmt.Println("Terima kasih sudah menggunakan aplikasi ini")
		return
	}
}
