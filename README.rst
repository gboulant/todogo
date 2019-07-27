The todogo application is a todo list manager written with the go
langage. The project was initiated as a training exercice to learn the
basic features of the go langage. The resulting application can be
used of course but could not fit your own way to manage todo list.

The standard use cases are to create and list a personal todo list.

For creating a new task to do::

  $ todo new -t "Write the documentation of todogo"
   1 [2019-07-27]  todo: Write the documentation of todogo

The todo task is created with usage index 1. You may complete your
todo list::
  
  $ todo new -t "Write the unit tests of todogo"
   2 [2019-07-27]  todo: Write the unit tests of todogo
  $ todo new -t "Create a beautiful web site for todogo"
   3 [2019-07-27]  todo: Create a beautiful web site for todogo

Each task is assigned a usage index (1,2,3) to refer to them in the
folowing use cases.   

Then you can have a look on the whole todo list::

  $ todo list
   1 [2019-07-27]  todo: Write the documentation of todogo
   2 [2019-07-27]  todo: Write the unit tests of todogo
   3 [2019-07-27]  todo: Create a beautiful web site for todogo

This week, you plan to work on the documentation (task 1) and the web
site (task 3), then you can star these tasks in putting them on the
board::

  $ todo board -a 1,3
  $ Task of index 1 has been added on board
  $ Task of index 3 has been added on board

And see the task on the board::

  $ todo board
   1 [2019-07-27]  todo: Write the documentation of todogo
   3 [2019-07-27]  todo: Create a beautiful web site for todogo

You start by writing some documentation and want to point that the
task is starting, then you specify that you jump to the next status of
this task, the status *doing*::

  $ todo status -n 1
   1 [2019-07-27] doing: Write the documentation of todogo

Then the board indicate::

  $ todo board
   1 [2019-07-27] doing: Write the documentation of todogo
   3 [2019-07-27]  todo: Create a beautiful web site for todogo

Note that todogo defines only three possible status: todo (the task is
registered and is waiting to be done), doing (the task is starting and
is in progress), done (the task is terminated)::

  $ todo status -n 1
   1 [2019-07-27]  done: Write the documentation of todogo

You can now get rid of this task from the board at least::

  $ todo board -r 1
  Task of index 1 has been removed from board

The task is always in the todo list (with status done), but no longer
on the board::

  $ todo list
   1 [2019-07-27]  done: Write the documentation of todogo
   2 [2019-07-27]  todo: Write the unit tests of todogo
   3 [2019-07-27]  todo: Create a beautiful web site for todogo
  $ todo board
   3 [2019-07-27]  todo: Create a beautiful web site for todogo

If you finish a lot of tasks, they could accumulate in your todo list,
with increasing indeces. A good practice is then to archive the done
tasks::

  $ todo archive -a 1
  Task 1 moved to the archive with a new usage index: 201907271605190109

Then the todo list is now::

  $ todo list
   2 [2019-07-27]  todo: Write the unit tests of todogo
   3 [2019-07-27]  todo: Create a beautiful web site for todogo

