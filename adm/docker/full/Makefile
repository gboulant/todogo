all: test

build:
	./dkbuild.sh

test: build
	./dktodo.sh config -n test
	./dktodo.sh new -t 'Create documentatin for todogo'
	./dktodo.sh list

clean:
	rm -f *~ dockerfiles/*~
	./dkclean.sh
