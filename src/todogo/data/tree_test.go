package data

import (
	"fmt"
	"testing"
)

func createTreeTaskArray() TaskArray {
	tasks := TaskArray{
		Task{UIndex: 30, Description: "A", ParentID: NoUID},
		Task{UIndex: 31, Description: "B", ParentID: NoUID},
		Task{UIndex: 32, Description: "C", ParentID: NoUID},
		Task{UIndex: 4, Description: "A.1", ParentID: NoUID},
		Task{UIndex: 5, Description: "A.2", ParentID: NoUID},
		Task{UIndex: 6, Description: "A.3", ParentID: NoUID},
		Task{UIndex: 7, Description: "B.1", ParentID: NoUID},
		Task{UIndex: 8, Description: "B.2", ParentID: NoUID},
		Task{UIndex: 9, Description: "C.1", ParentID: NoUID},
		Task{UIndex: 10, Description: "C.2", ParentID: NoUID},
		Task{UIndex: 11, Description: "B.2.1", ParentID: NoUID},
		Task{UIndex: 12, Description: "B.2.2", ParentID: NoUID},
		Task{UIndex: 13, Description: "B.2.2.1", ParentID: NoUID},
		Task{UIndex: 14, Description: "D", ParentID: NoUID},
	}

	taskFromText := func(description string) *Task {
		for i := 0; i < len(tasks); i++ {
			task := &tasks[i]
			if task.Description == description {
				return task
			}
		}
		return nil
	}

	// Set the parent dependency
	taskFromText("A.1").ParentID = taskFromText("A").UIndex
	taskFromText("A.2").ParentID = taskFromText("A").UIndex
	taskFromText("A.3").ParentID = taskFromText("A").UIndex
	taskFromText("B.1").ParentID = taskFromText("B").UIndex
	taskFromText("B.2").ParentID = taskFromText("B").UIndex
	taskFromText("B.2.1").ParentID = taskFromText("B.2").UIndex
	taskFromText("B.2.2").ParentID = taskFromText("B.2").UIndex
	taskFromText("C.1").ParentID = taskFromText("C").UIndex
	taskFromText("C.2").ParentID = taskFromText("C").UIndex
	taskFromText("B.2.2.1").ParentID = taskFromText("B.2.2").UIndex

	return tasks
}

func TestTaskArrayAncestor(t *testing.T) {
	viewlog = false

	tasks := createTreeTaskArray()
	printlog(tasks.String())

	taskFromText := func(description string) *Task {
		for i := 0; i < len(tasks); i++ {
			task := &tasks[i]
			if task.Description == description {
				return task
			}
		}
		return nil
	}

	childText := "C.2"
	parentText := "C.2"
	childID := taskFromText(childText).UIndex
	parentID := taskFromText(parentText).UIndex
	result := tasks.ancestor(childID, parentID)
	printlog(fmt.Sprintf("ancestor[c=%s, p=%s]: %v", childText, parentText, result))
	expect := false
	if result != expect {
		t.Errorf("The result is %v (should be %v)", result, expect)
	}

	childText = "B.2.2.1"
	parentText = "C"
	childID = taskFromText(childText).UIndex
	parentID = taskFromText(parentText).UIndex
	result = tasks.ancestor(childID, parentID)
	printlog(fmt.Sprintf("ancestor[c=%s, p=%s]: %v", childText, parentText, result))
	expect = false
	if result != expect {
		t.Errorf("The result is %v (should be %v)", result, expect)
	}

	childText = "B.2.2.1"
	parentText = "B"
	childID = taskFromText(childText).UIndex
	parentID = taskFromText(parentText).UIndex
	result = tasks.ancestor(childID, parentID)
	printlog(fmt.Sprintf("ancestor[c=%s, p=%s]: %v", childText, parentText, result))
	expect = true
	if result != expect {
		t.Errorf("The result is %v (should be %v)", result, expect)
	}
}

func TestTaskArrayTreeString(t *testing.T) {
	viewlog = false

	tasks := createTreeTaskArray()
	printlog(tasks.String())

	printlog("====================================")
	printlog("The resulting tree string is:")
	printlog("====================================")
	printlog(TreeString(tasks))

	taskFromText := func(description string) *Task {
		for i := 0; i < len(tasks); i++ {
			task := &tasks[i]
			if task.Description == description {
				return task
			}
		}
		return nil
	}

	// Simulate the cycle dependency
	taskFromText("B.2").ParentID = taskFromText("B.2.2.1").UIndex
	printlog("====================================")
	printlog("The resulting tree string is:")
	printlog("====================================")
	printlog(TreeString(tasks))

	// Restore the dependency
	taskFromText("B.2").ParentID = taskFromText("B").UIndex
	printlog("====================================")
	printlog("The resulting tree string is:")
	printlog("====================================")
	printlog(TreeString(tasks))

	// Change of branch
	taskFromText("B.2").ParentID = taskFromText("A.3").UIndex
	printlog("====================================")
	printlog("The resulting tree string is:")
	printlog("====================================")
	printlog(TreeString(tasks))

	taskFromText("B").ParentID = taskFromText("C.2").UIndex
	taskFromText("A").ParentID = taskFromText("B.1").UIndex
	taskFromText("C").ParentID = taskFromText("D").UIndex
	printlog("====================================")
	printlog(TreeString(tasks))

	viewlog = false

}
