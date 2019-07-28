package main

import (
	"fmt"
	"todogo/data"
)

func TestTask() {
	journal := data.CreateTestJournal()
	fmt.Println(journal.String())

	journalpath := "/tmp/todojournal.json"
	err := journal.SaveTo(journalpath)
	if err != nil {
		fmt.Println(err)
	}

	var anotherJournal data.TaskJournal
	anotherJournal.Load(journalpath)
	fmt.Println(anotherJournal.String())
	for i := 0; i < len(anotherJournal.TaskList); i++ {
		gindexInit := journal.TaskList[i].GIndex
		gindexRead := anotherJournal.TaskList[i].GIndex
		if gindexRead != gindexRead {
			err := fmt.Errorf("GIndex is %d (should be %d)", gindexRead, gindexInit)
			fmt.Println(err)
		}
	}
	anotherJournal.Save()
}

func main() {
	TestTask()
}
