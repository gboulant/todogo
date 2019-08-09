package core

// -----------------------------------------------------------------
// Definition of the jsonable interface

// Jsonable specifies the interface to be implemented by any structure to be
// serialized to a json file.
type Jsonable interface {
	Load(filepath string) error
	File() string
	Save() error
	SaveTo(filepath string) error
}

// Constant parameters for pretty presentation of json body
const (
	JSONPrefix string = ""
	JSONIndent string = "    "
)

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
