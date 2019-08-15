The todogo application is yet another todo list manager written with
the go langage. The project was initiated as a training exercice to
learn the basic features of the go langage. The resulting application
can be used of course but could not fit your own way to manage todo
list.

Todogo is a command line application, that consists in a single
executable program (todo) interacting with a local database made of
json files. It comes with no dependency (only the standard go packages
are required).

A typical output looks like:

```shell
$ todo list

 1 [2019-Aug-15] ▶ : Write the documentation of todogo
 2 [2019-Aug-15] ▶ : Write the unit tests of todogo
 3 [2019-Aug-15] ○ : Create a beautiful web site for todogo
 4 [2019-Aug-15] ● : Develop a tree representation of parent-children tasks
 5 [2019-Aug-15] ○ : Push a clone of the repository on github
 6 [2019-Aug-15] ● : Create a dockerfile of the todogo application

Legend: ○ todo  ▶ doing  ● done

```

Want to discover? Have a look to the user documentation:

* [Userguide](doc/userguide.rst)

For a rapid overview, the usage message below lists the existing features:

``` shell
$ todo -h
usage: todo <command> [<options>] [<arguments>]

With <command> in:

* new       : Create a new task
* list      : Print the list of tasks
* status    : Change the status of tasks
* board     : Append/Remove tasks on/from the board
* note      : Edit/View the note associated to a task
* child     : Make tasks be children of a parent task
* delete    : Delete tasks (definitely or in archive)
* archive   : Archive/Restore tasks
* config    : Manage de configuration and the contexts
* info      : Print detailled information on a task

For a description of possible options, try: todo <command> -h
```
