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

// --------------------------------------------------------------------------------------
// Functions for pretty printing of tasks

var renderingFunction func(s string, status TaskStatus) string
var renderingMap map[TaskStatus]string

func initRenderingTools() {
	config, _ := conf.GetConfig() // unused to test the err, we can not arrive here in case of config error
	if config.Parameters.PrettyPrint {
		renderingMap = taskStatusPretty
	} else {
		renderingMap = taskStatusSymbol
	}

	if config.Parameters.WithColor {
		renderingFunction = func(s string, status TaskStatus) string {
			return core.ColorString(s, taskStatusColors[status])
		}
	} else {
		renderingFunction = func(s string, status TaskStatus) string {
			return s
		}
	}
}

func (status TaskStatus) String() string {
	if renderingFunction == nil {
		initRenderingTools()
	}
	return renderingFunction(renderingMap[status], status)
}

func (status TaskStatus) legend() string {
	if renderingFunction == nil {
		initRenderingTools()
	}
	legend := fmt.Sprintf("%s %s", renderingMap[status], status.Label())
	return renderingFunction(legend, status)
}
