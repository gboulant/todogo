package data

import (
	"path/filepath"
	"testing"
	"todogo/core"
)

func TestCreateTask(t *testing.T) {
	var reference int64 = 1563794667
	task := Task{UIndex: 1, Timestamp: reference, Description: "Une tâche de test"}
	if task.Timestamp != reference {
		t.Errorf("timestamp is %d (should be %d)", task.Timestamp, reference)
	}
}

func TestJsonParsing(t *testing.T) {
	var ta TaskArray
	jsonfilepath, err := filepath.Abs("tasklist.json")
	if err == nil {
		err = core.Load(jsonfilepath, &ta)
	}

	if err != nil {
		t.Error(err)
	}

	err = core.Save("out.tasklist.json", &ta)
	if err != nil {
		t.Error(err)
	}
}
