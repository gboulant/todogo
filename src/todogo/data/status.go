package data

import (
	"errors"
	"fmt"
	"todogo/conf"
	"todogo/core"
)

// TaskStatus is an index of the step of completion of a task
type TaskStatus int

// Enumeration of possible TaskStatus
const (
	StatusTodo  TaskStatus = 0
	StatusDoing TaskStatus = 1
	StatusDone  TaskStatus = 2
	StatusStart TaskStatus = StatusTodo
	StatusEnd   TaskStatus = StatusDone
)

var taskStatusLabels = map[TaskStatus]string{
	StatusTodo:  "todo",
	StatusDoing: "doing",
	StatusDone:  "done",
}

var taskStatusColors = map[TaskStatus]core.ColorIndex{
	StatusTodo:  core.ColorGreen,
	StatusDoing: core.ColorOrange,
	StatusDone:  core.ColorBlue,
}

var taskStatusSymbol = map[TaskStatus]string{
	StatusTodo:  "o",
	StatusDoing: ">",
	StatusDone:  "x",
}

var taskStatusPretty = map[TaskStatus]string{
	StatusTodo:  core.PrettyDiskVoid,
	StatusDoing: core.PrettyTriangleRight,
	StatusDone:  core.PrettyDiskFull,
}

// Label returns a string representation of this status
func (status TaskStatus) Label() string {
	return taskStatusLabels[status]
}

// Value sets the status value from its string label
func (status *TaskStatus) Value(label string) error {
	for key, value := range taskStatusLabels {
		if label == value {
			*status = key
			return nil
		}
	}
	msg := fmt.Sprintf("ERR: the status %s is not defined (should be one of %v)",
		label, taskStatusLabels)
	return errors.New(msg)
}

func (status TaskStatus) String() string {
	if conf.PrettyPrint {
		return status.prettyString()
	} else {
		return status.symbolString()
	}
}

func (status TaskStatus) prettyString() string {
	if conf.WithColor {
		return core.ColorString(taskStatusPretty[status], taskStatusColors[status])
	}
	return taskStatusPretty[status]
}

func (status TaskStatus) symbolString() string {
	if conf.WithColor {
		return core.ColorString(taskStatusSymbol[status], taskStatusColors[status])
	} else {
		return taskStatusSymbol[status]
	}
}

func (status TaskStatus) legend() string {
	if conf.PrettyPrint {
		return status.prettyLegend()
	} else {
		return status.symbolLegend()
	}
}

func (status TaskStatus) prettyLegend() string {
	legend := fmt.Sprintf("%s %s", taskStatusPretty[status], status.Label())
	if conf.WithColor {
		return core.ColorString(legend, taskStatusColors[status])
	}
	return legend
}

func (status TaskStatus) symbolLegend() string {
	statusSymbol := status.symbolString()
	legend := fmt.Sprintf("%s %s", statusSymbol, status.Label())
	return legend
}

// Next makes the status change to its next state
func (status *TaskStatus) Next() error {
	if *status == StatusEnd {
		return errors.New("ERR: the status is already on the ending state")
	}
	*status++
	return nil
}

// Previous makes the status change to its previous state
func (status *TaskStatus) Previous() error {
	if *status == StatusStart {
		return errors.New("ERR: the status is already on the first state")
	}
	*status--
	return nil
}
