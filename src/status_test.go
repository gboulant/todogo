package todo

import (
	"testing"
)

func TestTaskStatus(t *testing.T) {
	var status TaskStatus
	status = StatusDoing
	//fmt.Printf("Status = %d (%s)\n", status, status.Label())
	if status.Label() != "doing" {
		t.Errorf("status label is %s (should be %s)", status.Label(), "doing")
	}

	status.Value("done")
	//fmt.Printf("Status = %d (%s)\n", status, status.Label())
	if status != StatusDone {
		t.Errorf("status is %v (should be %v)", status, StatusDone)
	}

	status = 5
	//fmt.Printf("Status = %d (%s)\n", status, status.Label())
	if status.Label() != "" {
		t.Errorf("status label is %s (should be %s)", status.Label(), "")
	}

	err := status.Value("closed")
	if err == nil {
		t.Error("closed is not a possible value for a status")
	}
}
