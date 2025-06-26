package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

type ExampleModel struct {
	gorm.Model
	Name string
}

func task(db *gorm.DB) {
	now := time.Now()
	log.Println("Task is being run... ", now)

	newRecord := ExampleModel{Name: fmt.Sprintf("Record at %s", now.Format(time.RFC3339))}
	db.Create(&newRecord)
}
