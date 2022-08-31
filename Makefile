
all: test

test:
	@make -C src test

clean:
	@make -C src clean
	@make -C doc/talks clean

dochtml:
	make -C doc/talks

# ------------------------------------------
# For installing the todo command and additionnal helper scripts
# (completion, git, synchro, etc) outside of the go bin directory, use
# the following targets (after setting the PREFIX value to specify the
# install root directory)

PREFIX ?= /usr/local
${PREFIX}/%:
	@mkdir -p $@

install: build ${PREFIX}/bin ${PREFIX}/etc
	install ./src/cmds/todo/todo ${PREFIX}/bin/.
	install ./adm/todo-git.sh ${PREFIX}/bin/.
	install ./adm/todo-sync.sh ${PREFIX}/bin/.
	install ./adm/todo-cfg.sh ${PREFIX}/bin/.
	install ./adm/todo-completion.sh ${PREFIX}/etc/.

build:
	@make -C src/cmds/todo build

