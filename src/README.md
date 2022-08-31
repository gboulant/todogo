This folder contains the todogo application source files.

The todogo application is composed of:

* A set of functions packaged as a go module of name
  `galuma.net/todo`. The source code is in this folder.
* An executable program of name `todo` that depends on the package
  `galuma.net/todo`. The source code is in the sub-folder `cmds/todo`.

For testing, compiling and installing the todo application, just use
the directives of the go toolchain:

```shell
$ go test       # test the package galuma.net/todo
$ cd cmds/todo  # change to the executable program directory
$ go build      # build the executable program for testing
$ go install    # install the executable program in the go bin directory
```

A Makefile is provided for additional technical actions.

You may also use `go doc` to get documentation of a function:

```shell
$ go doc todo.LoadBytes
package todo // import "galuma.net/todo"

func LoadBytes(fpath string) ([]byte, error)
    LoadBytes loads the content of the specified file and return it as an array
    of bytes.
	
```

And if you install the `godoc` tool, you can browse the documentation
online, by running this command at the root directory of the package
(this present directory): 

```shell
$ godoc -http=:8080
```

Then browse the URL `http://localhost:8080/pkg/galuma.net/todo`

**warning**: For installing the tool `godoc`, use the `go get`
instruction (see below) to be executed **outside** of the module
`galuma.net/todo`, otherwise it will just be added to the dependencies
of the module in the `go.mod` file.

```shell
$ go get -v golang.org/x/tools/cmd/godoc
```

**remark**: the `gotool` does not have an html generation feature (for
example to get the static html files for integration into a web
site). But it seems possible to sniff the html files from the running
`godoc` server using the `wget` command. Another possibility would be
to have a markdown export format with the `go doc` instruction (`go
doc`, with a space, and not the tool `godoc`). I think it does not
exist yet.
