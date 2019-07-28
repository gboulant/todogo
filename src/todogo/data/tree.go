package data

import "fmt"

type DataID uint64

// dataIDArray is a list of DataId
type dataIDArray []DataID

// treeMap is a map whose key is a data id and the value is the list of data id.
// This map is used to map the children data to their parents (the key is a
// parent data id and the value is the list of ID of its children)
type treeMap map[DataID]dataIDArray

// addChild adds a child in the list of children of the data of ID parentID
func (tree *treeMap) addChild(parentID DataID, childID DataID) error {
	if parentID == childID {
		return fmt.Errorf("ERR: a child can not be parent of itself (ID=%v)", childID)
	}
	_, exists := (*tree)[parentID]
	if !exists {
		(*tree)[parentID] = make(dataIDArray, 0)
	}
	(*tree)[parentID] = append((*tree)[parentID], childID)
	return nil
}

// initialize initializes the tree structure from the data array
func (tree *treeMap) initialize(tasks TaskArray) error {
	for i := 0; i < len(tasks); i++ {
		task := tasks[i]
		err := tree.addChild(DataID(task.ParentID), DataID(task.UIndex))
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

	const tabindent = "   "
	const tabchild = "└──"
	const tabstart = "* "
	startIndent := func(size int) string {
		s := ""
		for i := 0; i < size; i++ {
			s += " "
		}
		return s
	}(len(tabstart))

	// nodeString is the recursive function
	var nodeString func(dataID DataID, tab string) string
	nodeString = func(dataID DataID, tab string) string {
		idx := tasks.indexFromUID(uint64(dataID))
		s := fmt.Sprintf("%s%s\n", tab, tasks[idx].String())
		_, exists := tree[dataID]
		if !exists {
			return s
		}
		children := tree[dataID]
		if len(children) == 0 {
			return s
		}
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
	noParentDataIDs := tree[noUID]
	for k := 0; k < len(noParentDataIDs); k++ {
		stree += nodeString(noParentDataIDs[k], tabstart)
	}
	return stree
}
