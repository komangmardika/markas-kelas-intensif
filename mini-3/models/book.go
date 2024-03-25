package models

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

// point 7
type Book struct {
	Model
	ISBN    string `json:"isbn"`
	Penulis string `json:"penulis"`
	Tahun   uint   `json:"tahun"`
	Judul   string `json:"judul"`
	Gambar  string `json:"gambar"`
	Stok    uint   `json:"stok"`
}

func (book *Book) Update(db *gorm.DB) error {
	return db.Model(Book{}).Select("isbn", "penulis", "tahun", "judul", "gambar", "stok").
		Where("id = ?", book.Model.ID).
		Updates(map[string]interface{}{
			"isbn":    book.ISBN,
			"penulis": book.Penulis,
			"tahun":   book.Tahun,
			"judul":   book.Judul,
			"gambar":  book.Gambar,
			"stok":    book.Stok,
		}).
		Error
}

func (book *Book) SoftDeleteByIsbn(db *gorm.DB) error {
	return db.Model(Book{}).
		Where("isbn = ?", book.ISBN).
		Update("deleted_at", time.Now()).
		Error
}

func (book *Book) Create(db *gorm.DB) error {
	return db.Model(Book{}).
		Create(&book).Error
}

func (book *Book) CreateBulk(db *gorm.DB, books []Book) error {
	return db.Create(books).Error
}

func (book *Book) Upsert(db *gorm.DB) error {
	var b Book
	result := db.Model(Book{}).Where("id = ?", book.Model.ID).Take(&b)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		result := db.Create(book)
		if result.Error != nil {
			return result.Error
		}
	} else {
		result := db.Model(&b).Updates(book)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func (book *Book) GetById(db *gorm.DB) (Book, error) {
	return Book{}, db.Model(Book{}).
		Where("id = ?", book.Model.ID).
		Take(&book).
		Error
}

func (book *Book) GetByIsbn(db *gorm.DB) (Book, error) {
	return Book{}, db.Model(Book{}).
		Where("isbn = ?", book.ISBN).
		Take(&book).
		Error
}

func (book *Book) GetAll(db *gorm.DB) ([]Book, error) {
	var books []Book
	return books, db.Where("deleted_at IS NULL").Find(&books).
		Error
}

func (book *Book) DeleteById(db *gorm.DB) error {
	return db.Where("id = ?", book.Model.ID).
		Delete(&book).
		Error
}
