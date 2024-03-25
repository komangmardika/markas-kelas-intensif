package app

import (
	"encoding/csv"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"mini_3/configs"
	"mini_3/helpers"
	"mini_3/models"
	"os"
	"strconv"
	"sync"
)

// point 4
func ShowAll() {

	book := models.Book{}
	books, err := book.GetAll(configs.Mysql.DB)

	if err != nil {
		fmt.Println("terjadi kesalahan")
	} else {
		if len(books) == 0 {
			fmt.Println("belum ada buku, silakan tambah")
		} else {
			for i := 0; i < len(books); i++ {
				fmt.Println("No.", i+1)

				fmt.Printf(`
	ID: %d
	ISBN: %s
	Judul: %s
	Penulis: %s
	File Gambar: %s
	Tahun Publikasi: %d
	Stok: %d
`,
					books[i].Model.ID,
					books[i].ISBN,
					books[i].Judul,
					books[i].Penulis,
					books[i].Gambar,
					books[i].Tahun,
					books[i].Stok,
				)
			}
		}

	}
}

// point 1
func AddNewBook() {
	var bookJudul, bookIsbn, bookPenulis, bookGambar string
	var bookTahun, bookStok uint
	var books []models.Book

	for {
		for {
			bookIsbn = helpers.InputText("Masukkan ISBN buku baru: ")
			bookExists := isBookExistsByCode(bookIsbn)

			if !bookExists {
				break
			} else {
				fmt.Println("Buku dengan ISBN tersebut sudah ada, silakan mencoba lagi dengan kode ISBN lain")
				continue
			}
		}

		bookPenulis = helpers.InputText("Masukkan Penulis buku baru: ")
		bookJudul = helpers.InputText("Masukkan judul buku baru: ")
		bookGambar = helpers.InputText("Masukkan nama file gambar baru: ")
		bookTahun = helpers.InputUint("Masukkan tahun terbit buku baru: ")
		bookStok = helpers.InputUint("Masukkan stok buku baru: ")

		books = append(books, models.Book{
			ISBN:    bookIsbn,
			Penulis: bookPenulis,
			Judul:   bookJudul,
			Gambar:  bookGambar,
			Tahun:   bookTahun,
			Stok:    bookStok,
		})

		if helpers.KeyPress("Tekan y untuk menambah buku lainnya ") {
			continue
		} else {
			break
		}
	} //end for
	helpers.ClearScreen()

	if len(books) > 0 {
		fmt.Println("Menyimpan seluruh buku yang telah diinput")

		b := models.Book{}
		err := b.CreateBulk(configs.Mysql.DB, books)

		if err != nil {
			return
		} else {
			fmt.Println("buku-buku telah berhasil tersimpan")
		}

	} else {
		fmt.Println("tidak ada data yang dimasukkan, tidak menyimpan")
	}
}

// point 2
func EditBook() {
	bookIsbn := helpers.InputText("Masukkan ISBN buku yang akan disunting: ")
	b, res := findBookByIsbn(bookIsbn)
	if res {
		fmt.Println("ISBN ditemukan! : ", b.ISBN)
		b.Judul = helpers.InputText("Masukkan judul buku baru: ")
		b.Penulis = helpers.InputText("Masukkan penulis buku baru: ")
		b.Tahun = helpers.InputUint("Masukkan tahun terbit buku baru: ")
		b.Stok = helpers.InputUint("Masukkan stok buku baru: ")
		b.Gambar = helpers.InputText("Masukkan file gambar baru: ")

		err := b.Update(configs.Mysql.DB)
		if err != nil {
			return
		} else {
			fmt.Println("buku berhasil dimuktahirkan")
		}

	} else {
		fmt.Printf("buku dengan ISBN %s tidak ditemukan \n", bookIsbn)
		return
	}

	helpers.ClearScreen()
}

// point 3
func DeleteBook() {
	bookIsbn := helpers.InputText("Masukkan ISBN buku yang akan dihapus: ")
	exists := isBookExistsByCode(bookIsbn)
	if exists {
		book := models.Book{ISBN: bookIsbn}
		err := book.SoftDeleteByIsbn(configs.Mysql.DB)
		if err != nil {
			fmt.Println("Galat:", err)
			return
		}
	} else {
		fmt.Printf("buku dengan ISBN %s tidak ditemukan \n", bookIsbn)
		return
	}
	fmt.Println("Buku berhasil dihapus")
}

// point 5
func ImportBook() {
	fmt.Println("Sedang bekerja, mohon sabar menunggu ...")
	// buka file csv dan jangan lupa tutup
	file, err := os.Open("seeders/sample_books.csv")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("terjadi galat")
			panic(err)
		}
		fmt.Println("selesai memindahkan")
	}(file)

	// persiapan goroutine
	rowChannel := make(chan []string)
	bookChannel := make(chan models.Book)
	done := make(chan bool)
	var wg sync.WaitGroup

	// baca row demi row dan masukkan ke dalam channel
	go func() {
		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			panic(err)
		}

		for _, record := range records {
			rowChannel <- record
		}
		close(rowChannel)
	}()

	// 3 worker untuk didaftarkan ke dalam slice book
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for row := range rowChannel {
				book := parseRowToBook(row)
				bookChannel <- book
			}
		}()
	}

	// goroutine selesai
	go func() {
		wg.Wait()
		close(bookChannel)
		done <- true
	}()

	// Process book yang masuk ke dalam channel book untuk disimpan di slice books
	for book := range bookChannel {
		if err := book.Upsert(configs.Mysql.DB); err != nil {
			fmt.Printf("Error upserting book: %v\n", err)
		}
	}

	<-done

}
func parseRowToBook(row []string) models.Book {
	// Convert Tahun and Stok to uint
	tahun, _ := strconv.ParseUint(row[3], 10, 64)
	stok, _ := strconv.ParseUint(row[6], 10, 64)
	id, _ := strconv.ParseUint(row[0], 10, 64)
	if row[0] == "1" {
		println(id, row[0], uint(id))
	}
	return models.Book{
		Model:   models.Model{ID: uint(id)},
		ISBN:    row[1],
		Penulis: row[2],
		Tahun:   uint(tahun),
		Judul:   row[4],
		Gambar:  row[5],
		Stok:    uint(stok),
	}

}

// point 12
func PrintBook() {
	if helpers.KeyPress("Tekan y untuk mencetak seluruh buku ") {
		printAllBook()
	} else {
		bookIsbn := helpers.InputText("Masukkan ISBN Buku yang akan dicetak: ")
		book, res := findBookByIsbn(bookIsbn)
		if !res {
			fmt.Println("kesalahan terjadi, buku tidak ditemukan")
			return
		}

		printSingleBook(book)

	}
}

func Menu() {
	var menuChosen string

	helpers.ClearScreen()

	fmt.Printf(`
Aplikasi perpustakan sederhana 
	1. tampilkan semua buku
	2. tambah buku
	3. hapus sebuah buku
	4. sunting sebuah buku
	5. cetak buku
	6. impor csv
	7. Quit
`)
	menuChosen = helpers.InputText("Masukkan pilihan anda(1/2/3/4/5/6/7): ")

	switch menuChosen {
	case "1":
		ShowAll()
	case "2":
		AddNewBook()
	case "3":
		DeleteBook()
	case "4":
		EditBook()
	case "5":
		PrintBook()
	case "6":
		ImportBook()
	case "7":
		fmt.Println("Terima kasih telah menggunakan Aplikasi perpustakaan sederhana")
		os.Exit(0)
	default:
		fmt.Print("Kamu tidak memasukkan pilihan yang tepat, ")
	}

	confirmation()
}

func confirmation() {
	if helpers.KeyPress("Tekan y untuk kembali ke pilihan ") {
		Menu()
	}
}

func isBookExistsByCode(isbn string) bool {

	book := models.Book{ISBN: isbn}
	_, err := book.GetByIsbn(configs.Mysql.DB)
	if err != nil {
		return false
	}
	return true
}
func findBookByIsbn(isbn string) (models.Book, bool) {

	book := models.Book{ISBN: isbn}
	_, err := book.GetByIsbn(configs.Mysql.DB)
	if err != nil {
		return models.Book{}, false
	}
	return book, true
}

func printSingleBook(book models.Book) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "", 12)

	writeHeadingToPdf(pdf)
	writeBookToPDF(pdf, book)

	err := pdf.OutputFileAndClose(fmt.Sprintf("pdf/%s.pdf", book.ISBN))
	if err != nil {
		fmt.Println("Terjadi kesalahan ketika menyimpan pdf:", err)
		return
	}
}
func printAllBook() {

	b := models.Book{}
	books, err := b.GetAll(configs.Mysql.DB)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "", 10)
	writeHeadingToPdf(pdf)

	bookChan := make(chan models.Book)
	done := make(chan bool)

	go func() {
		for _, book := range books {
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

	go func() {
		wg.Wait()
		close(done)
	}()

	<-done

	err = pdf.OutputFileAndClose("pdf/all-books.pdf")
	if err != nil {
		fmt.Println("Error saving PDF:", err)
		return
	}
}

func writeHeadingToPdf(pdf *gofpdf.Fpdf) {
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(20, 10, "ISBN", "1", 0, "", true, 0, "")
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(40, 10, "JUDUL", "1", 0, "", true, 0, "")
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(40, 10, "PENULIS", "1", 0, "", true, 0, "")
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(40, 10, "GAMBAR", "1", 0, "", true, 0, "")
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(20, 10, "TAHUN", "1", 0, "R", true, 0, "")
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(20, 10, "STOK", "1", 0, "R", true, 0, "")
	pdf.Ln(-1)
}
func writeBookToPDF(pdf *gofpdf.Fpdf, book models.Book) {

	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(20, 10, book.ISBN, "1", 0, "", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(40, 10, book.Judul, "1", 0, "", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(40, 10, book.Penulis, "1", 0, "", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(40, 10, book.Gambar, "1", 0, "", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(20, 10, fmt.Sprintf("%d", book.Tahun), "1", 0, "R", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(20, 10, fmt.Sprintf("%d", book.Stok), "1", 0, "R", false, 0, "")
	pdf.Ln(-1)
}
