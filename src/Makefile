all: test

test:
	@go test -v

install:
	@make -C cmds/todo install

clean:
	@go clean
	@make -C cmds/todo clean
	@rm -f *~
