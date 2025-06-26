package main

import (
	"log"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Name        string `json:"name"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Price       uint   `json:"price"`
}

func createBook(db *gorm.DB, book *Book) error {
	result := db.Create(book)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func getBook(db *gorm.DB, id uint) (*Book, error) {
	var book Book
	result := db.First(&book, id)
	if result.Error != nil {
		return &book, result.Error
	}
	return &book, nil
}

func getBooks(db *gorm.DB) ([]Book, error) {
	var books []Book
	result := db.Find(&books)
	if result.Error != nil {
		return books, result.Error
	}
	return books, nil
}

func updateBook(db *gorm.DB, book *Book) error {
	result := db.Save(&book)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func deleteBook(db *gorm.DB, id uint) error {
	var book Book
	result := db.Delete(&book, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func searchBook(db *gorm.DB, bookName string) *Book {
	var book Book
	result := db.Where("name = ?", bookName).First(&book)
	if result.Error != nil {
		log.Fatalf("Error search book: %v", result.Error)
	}
	return &book
}

func searchBooks(db *gorm.DB, bookName string) []Book {
	var book []Book
	result := db.Where("name = ?", bookName).Order("price").Find(&book)
	if result.Error != nil {
		log.Fatalf("Error search book: %v", result.Error)
	}
	return book
}
