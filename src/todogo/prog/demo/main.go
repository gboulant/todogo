package main

import (
	"fmt"
	"todogo/data"
)

func TestDatabase() {
	var db data.Database
	err := db.Init(data.JournalPath)
	if err != nil {
		fmt.Println(err)
	}
	var t data.Task
	t = db.New("Acheter le pain")
	fmt.Printf("New task: %#v\n", t)
	t = db.New("Aller chercher les enfants")
	fmt.Printf("New task: %#v\n", t)
	t = db.New("Acheter du lait")
	fmt.Printf("New task: %#v\n", t)
	db.Commit()
}

func main() {
	//TestDatabase()
	//TestTaskStatus()
}
