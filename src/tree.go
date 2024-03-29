package todo

import (
	"fmt"
)

// treeMap is a map whose key is a data id and the value is the list of data id.
// This map is used to map the children data to their parents (the key is a
// parent data id and the value is the list of ID of its children)
type treeMap map[TaskID]TaskIDArray

// addChild adds a child in the list of children of the data of ID parentID
func (tree *treeMap) addChild(parentID TaskID, childID TaskID) error {
	if parentID == childID {
		return fmt.Errorf("ERR: a child can not be parent of itself (ID=%v)", childID)
	}
	_, exists := (*tree)[parentID]
	if !exists {
		(*tree)[parentID] = make(TaskIDArray, 0)
	}
	(*tree)[parentID] = append((*tree)[parentID], childID)
	return nil
}

// initialize initializes the tree structure from the data array
func (tree *treeMap) initialize(tasks TaskArray) error {
	for i := 0; i < len(tasks); i++ {
		task := tasks[i]
		err := tree.addChild(TaskID(task.ParentID), TaskID(task.UIndex))
		if err != nil {
			return err
		}
	}
	return nil
}

// TreeString returns a tree representation of the dataArray
func TreeString(tasks TaskArray) string {
	// Create the children tree
	tree := make(treeMap, 0)
	tree.initialize(tasks)

	// The principle is to iterate on dataArray root elements (element with no
	// parent) and then create a string reresentation of the children tree using
	// a recurcive function nodeString

	const groupsep = "\n" // group separator (a group is a task that has no parent and has children)
	const tabindent = "   "
	const tabchild = " └─"
	//const tabchild = " └" + core.PrettyArrowRight
	const tabstart = ""
	startIndent := func(size int) string {
		s := ""
		for i := 0; i < size; i++ {
			s += " "
		}
		return s
	}(len(tabstart))

	// nodeString is the recursive function
	var nodeString func(taskID TaskID, tab string) string
	nodeString = func(taskID TaskID, tab string) string {
		idx := tasks.indexFromUID(taskID)
		s := fmt.Sprintf("%s%s\n", tab, tasks[idx].String())

		// If the task is a main task (i.e. a task with no parent, which
		// can be determine by testing the current tabulation), then we
		// add a groupsep for a better interline between groups of tasks
		if tab == tabstart {
			s = groupsep + s
		}

		// Is there children tasks?
		_, exists := tree[taskID]
		if !exists {
			// This task has no child => stop the recurcive loop
			return s
		}
		children := tree[taskID]
		if len(children) == 0 {
			return s
		}
		// It is a project task (i.e. a task with no parent AND with
		// children tasks), then the tabulation of children should start
		// with a startIndent. Once more, it can be determine by testing
		// the current tabulation
		if tab == tabstart {
			tab = startIndent + tabchild
		} else {
			tab = tabindent + tab
		}

		for i := 0; i < len(children); i++ {
			s += nodeString(children[i], tab)
		}
		return s
	}

	stree := ""
	noParentTaskIDs := tree[NoUID]
	for k := 0; k < len(noParentTaskIDs); k++ {
		stree += nodeString(noParentTaskIDs[k], tabstart)
	}
	return stree
}
