package data

import (
	"fmt"
	"testing"
)

var viewlog = false

func printlog(msg string) {
	if viewlog {
		fmt.Println(msg)
	}
}

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
	printlog(tasks.String())

	index := tasks.indexFromUID(4)
	if index != 0 {
		t.Errorf("task index is %d (should be 0)", index)
	}
	tasks.SortByUID()
	printlog(tasks.String())
	index = tasks.indexFromUID(4)
	if index != 3 {
		t.Errorf("task index is %d (should be 3)", index)
	}
	tasks.SortByGID()
	printlog(tasks.String())
	tasks.SortByTimestamp()
	printlog(tasks.String())
}

func TestTaskArrayEdit(t *testing.T) {

	tasks := TaskArray{
		CreateTestTask(1, "Write documentation for todogo"),
		CreateTestTask(2, "Create unit test for todogo"),
		CreateTestTask(3, "Add a function to print a tasks journal"),
		CreateTestTask(4, "Organize a code review of todogo"),
	}
	printlog(tasks.String())

	ptask, _ := tasks.GetTask(2)
	if ptask.Description != "Create unit test for todogo" {
		msg := fmt.Sprintf("Description is \"%s\" (should be \"%s\")", ptask.Description, "Create unit test for todogo")
		t.Error(msg)
	}

	ptask.Description = "toto"
	otherTaskPointer, _ := tasks.GetTask(2)
	if otherTaskPointer.Description != "toto" {
		msg := fmt.Sprintf("Description is \"%s\" (should be \"%s\")", otherTaskPointer.Description, "toto")
		t.Error(msg)
	}
	printlog(tasks.String())

	tasks.Remove(tasks.indexFromUID(2))
	ptask, _ = tasks.GetTask(2)
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

	ptask, _ = tasks.GetTask(10)
	if ptask.Description != "Organize the documentation review" {
		msg := fmt.Sprintf("Description is \"%s\" (should be \"%s\")", ptask.Description, "Organize the documentation review")
		t.Error(msg)
	}
	tasks.SortByUID()
	printlog(tasks.String())
}

func TestTaskArrayFreeUID(t *testing.T) {

	tasks := TaskArray{
		CreateTestTask(1, "Write documentation for todogo"),
		CreateTestTask(2, "Create unit test for todogo"),
		CreateTestTask(5, "Add a function to print a tasks journal"),
		CreateTestTask(4, "Organize a code review of todogo"),
	}
	printlog(tasks.String())

	resuid := tasks.getFreeUID()
	refuid := uint64(3)
	printlog(fmt.Sprintf("uid: %d", resuid))
	if resuid != refuid {
		msg := fmt.Sprintf("UID is %d (should be %d)", resuid, refuid)
		t.Error(msg)
	}

	task := CreateTestTask(3, "Organize the documentation review")
	tasks.Append(task)
	resuid = tasks.getFreeUID()
	refuid = uint64(6)
	printlog(fmt.Sprintf("uid: %d", resuid))
	if resuid != refuid {
		msg := fmt.Sprintf("UID is %d (should be %d)", resuid, refuid)
		t.Error(msg)
	}

	printlog(tasks.String())

}

func TestTaskArrayFilter(t *testing.T) {
	viewlog = false

	tasks := TaskArray{
		CreateTestTask(1, "Write documentation for todogo"),
		CreateTestTask(2, "Create unit test for todogo"),
		CreateTestTask(5, "Add a function to print a tasks journal"),
		CreateTestTask(4, "Organize a code review of todogo"),
		CreateTestTask(3, "Create documentation review fro todogo"),
	}

	ptask, _ := tasks.GetTask(4)
	ptask.OnBoard = true
	tasksOnBoard := tasks.GetTasksWithFilter(TaskFilterOnBoard)
	for i := 0; i < len(tasksOnBoard); i++ {
		printlog(tasksOnBoard[i].String())
	}
	if len(tasksOnBoard) != 1 {
		msg := fmt.Sprintf("Nb tasks on board is %d (should be %d)", len(tasksOnBoard), 1)
		t.Error(msg)
	}
	ptask = tasksOnBoard[0]
	resuid := ptask.UIndex
	refuid := uint64(4)
	if resuid != refuid {
		msg := fmt.Sprintf("UID is %d (should be %d)", resuid, refuid)
		t.Error(msg)
	}

	viewlog = false
}

func TestTaskJournalEdit(t *testing.T) {

	journal := CreateTestJournal()

	ptask := journal.New("Setup the automatic daily test procedure")
	resuid := ptask.UIndex
	refuid := uint64(5)
	if resuid != refuid {
		msg := fmt.Sprintf("UID is %d (should be %d)", resuid, refuid)
		t.Error(msg)
	}
	printlog(journal.String())
}

func TestTaskJournalIO(t *testing.T) {
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
