package models_test

import (
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"mini_3/app"
	"mini_3/configs"
	"mini_3/models"
	"testing"
)

func Init() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("env not found, using system env")
	}
	configs.OpenDB(false)
}

func TestBook_Create(t *testing.T) {
	Init()
	book := models.Book{
		ISBN:    app.GenerateISBN(),
		Judul:   faker.Sentence(),
		Penulis: faker.FirstName() + " " + faker.LastName(),
		Tahun:   uint(rand.Intn(30) + 1990),
		Stok:    uint(rand.Intn(30)),
		Gambar:  faker.Word() + ".png",
	}

	err := book.Create(configs.Mysql.DB)
	assert.Nil(t, err)
}

// point 10
func TestBook_CreateBulk(t *testing.T) {
	Init()
	var books []models.Book

	maxi := rand.Intn(5) + 1
	for i := 0; i < maxi; i++ {
		books = append(books, models.Book{
			ISBN:    app.GenerateISBN(),
			Judul:   faker.Sentence(),
			Penulis: faker.FirstName() + " " + faker.LastName(),
			Tahun:   uint(rand.Intn(30) + 1990),
			Stok:    uint(rand.Intn(30)),
			Gambar:  faker.Word() + ".png",
		})
	}

	book := models.Book{}
	err := book.CreateBulk(configs.Mysql.DB, books)
	assert.Nil(t, err)
}

func TestBook_Upsert(t *testing.T) {
	Init()

	var books []models.Book

	maxi := rand.Intn(5) + 1
	for i := 0; i < maxi; i++ {
		books = append(books, models.Book{
			ISBN:    app.GenerateISBN(),
			Judul:   faker.Sentence(),
			Penulis: faker.FirstName() + " " + faker.LastName(),
			Tahun:   uint(rand.Intn(30) + 1990),
			Stok:    uint(rand.Intn(30)),
			Gambar:  faker.Word() + ".png",
		})
	}

	book := models.Book{}
	err := book.CreateBulk(configs.Mysql.DB, books)
	assert.Nil(t, err)

	entry := models.Book{
		Model:   models.Model{ID: 1},
		ISBN:    app.GenerateISBN(),
		Judul:   faker.Sentence(),
		Penulis: faker.FirstName() + " " + faker.LastName(),
		Tahun:   uint(rand.Intn(30) + 1990),
		Stok:    uint(rand.Intn(30)),
		Gambar:  faker.Word() + ".png",
	}

	err = entry.Upsert(configs.Mysql.DB)
	assert.Nil(t, err)
	assert.NotEqual(t, book.ISBN, entry.ISBN)
}

func TestBook_DeleteById(t *testing.T) {
	Init()
	book := models.Book{
		ISBN:    app.GenerateISBN(),
		Judul:   faker.Sentence(),
		Penulis: faker.FirstName() + " " + faker.LastName(),
		Tahun:   uint(rand.Intn(30) + 1990),
		Stok:    uint(rand.Intn(30)),
		Gambar:  faker.Word() + ".png",
	}

	err := book.Create(configs.Mysql.DB)
	assert.Nil(t, err)

	err = book.DeleteById(configs.Mysql.DB)
	assert.Nil(t, err)

}

// point 10
func TestBook_DeleteByIsbn(t *testing.T) {
	Init()
	book := models.Book{
		ISBN:    app.GenerateISBN(),
		Judul:   faker.Sentence(),
		Penulis: faker.FirstName() + " " + faker.LastName(),
		Tahun:   uint(rand.Intn(30) + 1990),
		Stok:    uint(rand.Intn(30)),
		Gambar:  faker.Word() + ".png",
	}

	err := book.Create(configs.Mysql.DB)
	assert.Nil(t, err)

	err = book.SoftDeleteByIsbn(configs.Mysql.DB)
	assert.Nil(t, err)
}

// point 10
func TestBook_GetAll(t *testing.T) {
	Init()
	book := models.Book{
		ISBN:    app.GenerateISBN(),
		Judul:   faker.Sentence(),
		Penulis: faker.FirstName() + " " + faker.LastName(),
		Tahun:   uint(rand.Intn(30) + 1990),
		Stok:    uint(rand.Intn(30)),
		Gambar:  faker.Word() + ".png",
	}

	err := book.Create(configs.Mysql.DB)

	books, err := book.GetAll(configs.Mysql.DB)

	assert.GreaterOrEqual(t, len(books), 1)
	assert.Nil(t, err)

	err2 := book.SoftDeleteByIsbn(configs.Mysql.DB)

	assert.Nil(t, err2)
	assert.GreaterOrEqual(t, len(books), 0)
}

func TestBook_GetById(t *testing.T) {
	Init()
	book := models.Book{
		ISBN:    app.GenerateISBN(),
		Judul:   faker.Sentence(),
		Penulis: faker.FirstName() + " " + faker.LastName(),
		Tahun:   uint(rand.Intn(30) + 1990),
		Stok:    uint(rand.Intn(30)),
		Gambar:  faker.Word() + ".png",
	}

	errNew := book.Create(configs.Mysql.DB)
	assert.Nil(t, errNew)

	b := models.Book{
		Model: models.Model{ID: book.Model.ID},
	}
	_, err := b.GetById(configs.Mysql.DB)
	assert.Nil(t, err)
	assert.Equal(t, b.ISBN, book.ISBN)
}

func TestBook_GetByIsbn(t *testing.T) {
	Init()
	book := models.Book{
		ISBN:    app.GenerateISBN(),
		Judul:   faker.Sentence(),
		Penulis: faker.FirstName() + " " + faker.LastName(),
		Tahun:   uint(rand.Intn(30) + 1990),
		Stok:    uint(rand.Intn(30)),
		Gambar:  faker.Word() + ".png",
	}

	errNew := book.Create(configs.Mysql.DB)
	assert.Nil(t, errNew)

	b := models.Book{
		ISBN: book.ISBN,
	}
	_, err := b.GetByIsbn(configs.Mysql.DB)
	assert.Nil(t, err)
	assert.Equal(t, b.Model.ID, book.Model.ID)
}

// point 10
func TestBook_Update(t *testing.T) {
	Init()
	oldIsbn := "123"
	oldJudul := "foo"
	oldPenulis := "bar"
	oldGambar := "baz"
	oldTahun := 1990
	oldStok := 1

	book := models.Book{
		ISBN:    oldIsbn,
		Judul:   oldJudul,
		Penulis: oldPenulis,
		Tahun:   uint(oldTahun),
		Stok:    uint(oldStok),
		Gambar:  oldGambar,
	}

	errNew := book.Create(configs.Mysql.DB)
	assert.Nil(t, errNew)

	book.ISBN = app.GenerateISBN()
	book.Judul = faker.Sentence()
	book.Penulis = faker.FirstName() + " " + faker.LastName()
	book.Tahun = uint(rand.Intn(30) + 1990)
	book.Stok = uint(rand.Intn(30))
	book.Gambar = faker.Word() + ".png"

	err := book.Update(configs.Mysql.DB)
	assert.Nil(t, err)
	assert.NotEqual(t, oldIsbn, book.ISBN)
}
