package data

import (
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
