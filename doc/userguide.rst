================
TODO: User guide
================


The standard use cases are to create and manage a personal todo
list. For creating a new task::

   $ todo new -t "Write the documentation of todogo"
    1 [2019-07-27]  todo: Write the documentation of todogo

The todo task is created with usage index 1. You may complete your
todo list::

   $ todo new -t "Write the unit tests of todogo"
    2 [2019-07-27]  todo: Write the unit tests of todogo

   $ todo new -t "Create a beautiful web site for todogo"
    3 [2019-07-27]  todo: Create a beautiful web site for todogo

Each task is assigned a usage index (1,2,3) to refer to them in the
folowing use cases. Then you can have a look on the whole todo list::

   $ todo list
    1 [2019-07-27]  todo: Write the documentation of todogo
    2 [2019-07-27]  todo: Write the unit tests of todogo
    3 [2019-07-27]  todo: Create a beautiful web site for todogo

This week you plan to work on the documentation (task 1) and the web
site (task 3), then you can star these tasks in putting them on the
board::

   $ todo board -a 1,3
   Task of index 1 has been added on board
   Task of index 3 has been added on board

And see the tasks on the board to focus on the actuality::

   $ todo board
    1 [2019-07-27]  todo: Write the documentation of todogo
    3 [2019-07-27]  todo: Create a beautiful web site for todogo

You start by writing some documentation and want to point that the
task is in progress, then you specify that you jump to the next status of
this task 1 (the status *doing*)::

   $ todo status -n 1
    1 [2019-07-27] doing: Write the documentation of todogo

Then the board indicates::

   $ todo board
    1 [2019-07-27] doing: Write the documentation of todogo
    3 [2019-07-27]  todo: Create a beautiful web site for todogo

Note that todogo defines three possible status: todo (the task is
registered and is waiting to be done), doing (the task is started and
is in progress), done (the task is achieved)::
   
   $ todo status -n 1
    1 [2019-07-27]  done: Write the documentation of todogo

You can now get rid of this task from the board::

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

If you register and then finish a lot of tasks, they could accumulate
in your todo list, with increasing indeces. A good practice is then to
archive the done tasks::

   $ todo archive -a 1
   Task 1 moved to the archive with a new usage index: 201907271605190109

Then the todo list is now::

   $ todo list
    2 [2019-07-27]  todo: Write the unit tests of todogo
    3 [2019-07-27]  todo: Create a beautiful web site for todogo

And the archive contains::

   $ todo archive
   201907271605190109 [2019-07-27]  done: Write the documentation of todogo

Note that when a task is moved to the archive, then its usage index is
modfied and set to its absolute index. When created, a task is
characterized by a usage index (the index seen by the user to
manipulate the task) and an absolute index (used by the program to
manage the tasks). The absolute index is unique and invariant ever for
a task all long of its life cycle. The usage index of a task is unique
and invariant as long as the task is in the journal. Once a task is
move from the journal to the archive, its usage index is realeased and
can be reused for a new task::

   $ todo new -t "Make it possible to have children tasks associated to a task"
    1 [2019-07-27]  todo: Make it possible to have children tasks associated to a task

As you can see, the usage index 1, previously attributed to the
documentation task (moved to the archive) has been recycled and
attributed to this newly created task::

   $ todo list
    1 [2019-07-27]  todo: Make it possible to have children tasks associated to a task
    2 [2019-07-27]  todo: Write the unit tests of todogo
    3 [2019-07-27]  todo: Create a beautiful web site for todogo

The reason of this index recycling is to avoid increasing indeces, at
least in the journal listing, so that you can refer to reasonably
short indeces when typing your command line. Even if there is no
maximum limit for indeces, the normal usage (i.e. if you achieve your
tasks and archive them when finished) is to play whith indeces between
1 (the starting index value) to 20 or 30.

The absolute indeces never changes whatever the location of the task
(journal or archive). It is defined as a concatenation of a date flag
YYYYMMDD and a sha1 of the task. For example, you have to manipulate
this absolute index to restore a task from the archive (for example in
the case where you forgot a part of the task)::

   $ todo archive -r 201907271605190109
   Task 201907271605190109 restored from archive with a new usage index: 4

As you can see, the task has been restored from the archive (where its
index was 201907271605190109, i.e. its absolute index) to the journal
with a new usage index 4 (of course the original index 1 has been
reassigned to another task and the first free usage index in the
journal is 4)::

   $ todo list
    1 [2019-07-27]  todo: Make it possible to have children tasks associated to a task
    2 [2019-07-27]  todo: Write the unit tests of todogo
    3 [2019-07-27]  todo: Create a beautiful web site for todogo
    4 [2019-07-27]  done: Write the documentation of todogo

The restored task is on status done, and it could be relevant to move
its status to the previous one in the sequence (the status "doing")::

   $ todo status -p 4
    4 [2019-07-27] doing: Write the documentation of todogo
