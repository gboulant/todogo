GOPATH=$(shell echo $$GOPATH:`pwd`)
PREFIX ?= /usr/local

all: build

demobuild:
	GOPATH=${GOPATH} go build todogo/prog/demo

demo: demobuild
	./demo

build:
	GOPATH=${GOPATH} go install todogo/prog/todo

install: build
	install ./bin/todo ${PREFIX}/bin/.
	install ./adm/todo-completion.sh ${PREFIX}/etc/.

test:
	@echo "=== Testing the package todogo/core ..."
	@GOPATH=${GOPATH} go test -v todogo/core
	@echo "=== Testing the package todogo/conf ..."
	@GOPATH=${GOPATH} go test -v todogo/conf
	@echo "=== Testing the package todogo/data ..."
	@GOPATH=${GOPATH} go test -v todogo/data

clean:
	rm -f ./todo ./demo
	find . -name "*~" -o -name "out.*" | xargs rm -f
	rm -rf ./pkg ./bin ./doc/api

edit:
	GOPATH=${GOPATH} code .


doc/api:
	@mkdir doc/api

docbuild: doc/api
	@GOPATH=${GOPATH} godoc -html todogo/core > ./doc/api/core.html
	@GOPATH=${GOPATH} godoc -html todogo/conf > ./doc/api/conf.html
	@GOPATH=${GOPATH} godoc -html todogo/data > ./doc/api/data.html

docview: docbuild
	@echo "Open the link http://localhost:6060/doc/api/data.html"
	@GOPATH=${GOPATH} godoc --http=:6060 -goroot=.
