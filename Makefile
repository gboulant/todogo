GOPATH=$(shell echo $$GOPATH:`pwd`)
PREFIX ?= ${HOME}/programs/local/bin

all: install

build:
	GOPATH=${GOPATH} go build todogo/prog/todo
	GOPATH=${GOPATH} go build todogo/prog/demo

install:
	GOPATH=${GOPATH} go install todogo/prog/todo

deploy: install
	install ./bin/todo ${PREFIX}

test:
	@echo "=== Testing the package todogo/core ..."
	@GOPATH=${GOPATH} go test -v todogo/core
	@echo "=== Testing the package todogo/data ..."
	@GOPATH=${GOPATH} go test -v todogo/data

clean:
	rm -f ./todo ./demo *~ out.*
	rm -f src/todogo/data/out.*
	rm -rf ./pkg ./bin ./doc/api

edit:
	GOPATH=${GOPATH} code .


doc/api:
	@mkdir doc/api

docbuild: doc/api
	@GOPATH=${GOPATH} godoc -html todogo/core > ./doc/api/core.html
	@GOPATH=${GOPATH} godoc -html todogo/data > ./doc/api/data.html

docview: docbuild
	@echo "Open the link http://localhost:6060/doc/api/data.html"
	@GOPATH=${GOPATH} godoc --http=:6060 -goroot=.
