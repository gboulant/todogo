The todogo application is yet another todo list manager written with
the go langage. The project was initiated as a training exercice to
learn the basic features of the go langage. The resulting application
can be used of course but could not fit your own way to manage todo
lists.

The todo program is a command line application, that consists in a
single executable program (todo) interacting with a local database
made of json files. It comes with no dependency (only the standard go
packages are required).

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

* [User guide](doc/talks/talk01.01.userguide.rst)
* [Getting started](doc/talks/talk01.02.gettingstarted.rst)
* [Basic design](doc/talks/talk01.03.basicdesign.rst)
