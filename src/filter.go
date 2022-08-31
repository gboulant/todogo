package todo

// TaskFilter defines a function that can be used to filter a list of task
// considering the return value (true or false) of the TaskFilter function.
type TaskFilter func(task Task) bool

// TaskFilterAll always returns true (no filter)
func TaskFilterAll(task Task) bool {
	return true
}

// TaskFilterOnBoard returns true if the task is on board
func TaskFilterOnBoard(task Task) bool {
	return task.OnBoard
}

// TaskFilterDone returns true if the task is done
func TaskFilterDone(task Task) bool {
	return (task.Status == StatusDone)
}

// TaskFilterDoing returns true if the task is doing
func TaskFilterDoing(task Task) bool {
	return (task.Status == StatusDoing)
}

// TaskFilterTodo returns true if the task is todo
func TaskFilterTodo(task Task) bool {
	return (task.Status == StatusTodo)
}
