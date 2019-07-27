package core

// -----------------------------------------------------------------
// Definition of the jsonable interface

type Jsonable interface {
	Load(filepath string) error
	Save(filepath string) error
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
	return container.Save(filepath)
}

// -----------------------------------------------------------------
// Definition of the stringable interface

type Stringable interface {
	String() string
	Println()
}

func String(data Stringable) string {
	return data.String()
}
func Println(data Stringable) {
	data.Println()
}
