package data

import (
	"testing"
)

func TestDateLabel(t *testing.T) {
	timestamp := int64(1563794667) // July 22, 2019
	timelabel := datelabel(timestamp)
	reference := "2019-Jul-22"
	if timelabel != reference {
		t.Errorf("timelabel is %s (should be %s)", timelabel, reference)
	}
}

func TestHashDate(t *testing.T) {
	timestamp := int64(1563794667) // July 22, 2019
	text := "1 [2019-07-22]: Le texte de la fiche"
	hash := hashdate(text, timestamp)
	reference := uint64(201907220287959902)
	if hash != reference {
		t.Errorf("hashdate is %d (should be %d)", hash, reference)
	}
}
