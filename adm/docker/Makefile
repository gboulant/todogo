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

minibuild:
	docker build -f dockerfiles/Dockerfile.minimalist -t todogomin .

minitest:
	docker run --rm -it todogomin make -C todogo test

minidemo:
	docker run --rm -it todogomin todo config -i
