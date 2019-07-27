package data

import (
	"testing"
)

func TestDatabase(t *testing.T) {
	var db Database
	err := db.Init("/tmp/journal.json")
	if err != nil {
		t.Error(err)
	}
	linit := len(db.taskmap)
	var task Task
	task = db.New("Acheter le pain")
	task = db.New("Aller chercher les enfants")
	task = db.New("Acheter du lait")
	if task.Description != "Acheter du lait" {
		t.Errorf("Description is %s (should be %s)", task.Description, "Acheter du lait")
	}

	reference := linit + 3
	if len(db.taskmap) != reference {
		t.Errorf("Nb task is %d (should be %d)", len(db.taskmap), reference)
	}
	db.Commit()
}
