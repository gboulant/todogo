package core

import "fmt"

// -----------------------------------------------------------------
// Definition of the jsonable interface

type Jsonable interface {
	Load(filepath string) error
	File() string
	Save() error
	SaveTo(filepath string) error
}

const JsonPrefix string = ""
const JsonIndent string = "    "

// Load reads a json file and maps the data to the jsonable container. The
// container should implements the Jsonable interface.
func Load(filepath string, container Jsonable) error {
	return container.Load(filepath)
}

// Save writes the data of the jsonable container into a file. The
// container should implements the Jsonable interface.
func Save(filepath string, container Jsonable) error {
	return container.SaveTo(filepath)
}

// -----------------------------------------------------------------
// Definition of the stringable interface

type Stringable interface {
	String() string
}

func String(data Stringable) string {
	return data.String()
}
func Println(data Stringable) {
	fmt.Println(data.String())
}
