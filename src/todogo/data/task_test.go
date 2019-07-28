package data

import (
	"fmt"
	"testing"
)

func TestTask(t *testing.T) {
	task := CreateTestTask(10, "Write documentation for todogo")

	status := task.Status
	if status != StatusStart {
		t.Errorf("Status is %d (should be %d)", status, StatusStart)
	}

	task.Status.Next()
	label := task.Status.Label()
	if label != "doing" {
		t.Errorf("Status label is %s (should be %s)", label, "doing")
	}
}

func TestTaskArraySort(t *testing.T) {
	tasks := TaskArray{
		CreateTestTask(4, "Write documentation for todogo"),
		CreateTestTask(3, "Create unit test for todogo"),
		CreateTestTask(2, "Add a function to print a tasks journal"),
		CreateTestTask(1, "Organize a code review of todogo"),
	}
	//fmt.Println(tasks.String())

	index := tasks.IndexFromUID(4)
	if index != 0 {
		t.Errorf("task index is %d (should be 0)", index)
	}
	tasks.SortByUID()
	//fmt.Println(tasks.String())
	index = tasks.IndexFromUID(4)
	if index != 3 {
		t.Errorf("task index is %d (should be 3)", index)
	}
	tasks.SortByGID()
	//fmt.Println(tasks.String())
	tasks.SortByTimestamp()
	//fmt.Println(tasks.String())
}

func TestTaskArrayEdit(t *testing.T) {

	log := func(msg string) {
		//fmt.Println(msg)
	}

	tasks := TaskArray{
		CreateTestTask(1, "Write documentation for todogo"),
		CreateTestTask(2, "Create unit test for todogo"),
		CreateTestTask(3, "Add a function to print a tasks journal"),
		CreateTestTask(4, "Organize a code review of todogo"),
	}
	log(tasks.String())

	ptask := tasks.GetTask(2)
	if ptask.Description != "Create unit test for todogo" {
		msg := fmt.Sprintf("Description is \"%s\" (should be \"%s\")", ptask.Description, "Create unit test for todogo")
		t.Error(msg)
	}

	ptask.Description = "toto"
	otherTaskPointer := tasks.GetTask(2)
	if otherTaskPointer.Description != "toto" {
		msg := fmt.Sprintf("Description is \"%s\" (should be \"%s\")", otherTaskPointer.Description, "toto")
		t.Error(msg)
	}
	log(tasks.String())

	tasks.Remove(tasks.IndexFromUID(2))
	ptask = tasks.GetTask(2)
	if ptask != nil {
		msg := fmt.Sprintf("Task %d should not exist", ptask.UIndex)
		t.Error(msg)
	}

	initlen := len(tasks)
	task := CreateTestTask(10, "Organize the documentation review")
	tasks.Append(task)
	if len(tasks) != initlen+1 {
		msg := fmt.Sprintf("Tasks size is %d (should be %d)", len(tasks), initlen+1)
		t.Error(msg)
	}

	ptask = tasks.GetTask(10)
	if ptask.Description != "Organize the documentation review" {
		msg := fmt.Sprintf("Description is \"%s\" (should be \"%s\")", ptask.Description, "Organize the documentation review")
		t.Error(msg)
	}
	tasks.SortByUID()
	log(tasks.String())
}

func TestTaskJournal(t *testing.T) {
	journal := CreateTestJournal()

	journalpath := "/tmp/todojournal.json"
	err := journal.SaveTo(journalpath)
	if err != nil {
		t.Error(err)
	}

	var anotherJournal TaskJournal
	anotherJournal.Load(journalpath)
	for i := 0; i < len(anotherJournal.TaskList); i++ {
		gindexInit := journal.TaskList[i].GIndex
		gindexRead := anotherJournal.TaskList[i].GIndex
		if gindexRead != gindexRead {
			t.Errorf("GIndex is %d (should be %d)", gindexRead, gindexInit)
		}
	}
	err = anotherJournal.Save()
	if err != nil {
		t.Error(err)
	}

}
