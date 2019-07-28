package main

import (
	"fmt"
	"todogo/data"
)

func TestTask() {
	tasks := data.TaskArray{
		data.CreateTestTask(4, "Write documentation for todogo"),
		data.CreateTestTask(3, "Create unit test for todogo"),
		data.CreateTestTask(2, "Add a function to print a tasks journal"),
		data.CreateTestTask(1, "Organize a code review of todogo"),
	}
	fmt.Println(tasks.String())

	index := tasks.IndexFromUID(4)
	if index != 0 {
		fmt.Printf("task index is %d (should be 0)", index)
	}
	tasks.SortByUID()
	fmt.Println(tasks.String())
	index = tasks.IndexFromUID(4)
	if index != 3 {
		fmt.Printf("task index is %d (should be 3)", index)
	}

	tasks.SortByGID()
	fmt.Println(tasks.String())

	tasks.SortByTimestamp()
	fmt.Println(tasks.String())

}

func main() {
	TestTask()
}
