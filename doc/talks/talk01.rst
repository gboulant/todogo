:Title: TODOGO - Quick start guide
:Author: Guillaume Boulant (gboulant@gmail.com)
:Date: Sept. 2019
:Description: Quick start guide for TODOGO
   
-------------

.. raw:: html

   <div align="center" style="padding-top: 20%; padding-left:20%; padding-right:20%">
   <h1 style="margin-left: 0; margin-right: 0">01 - Userguide</h1>
   <p style="margin-left: 0; margin-right: 0">
   </p>
   </div>

   
========
Overview
========

The todo program is a command line application created to manage a
personal todo list from a shell terminal:

.. code:: shell

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
   * config    : Manage de configuration

   For a description of possible options, try: todo <command> --help
   
========================
Creating tasks - ``new``
========================

For creating a new task:

.. code:: shell

   $ todo new -t "Write the documentation of todogo"
    1 [2019-Sep-07] ○ : Write the documentation of todogo

The todo task is created with usage index 1. You may complete your
todo list:

.. code:: shell

   $ todo new -t "Write the unit tests of todogo"
    2 [2019-Sep-07] ○ : Write the unit tests of todogo

   $ todo new -t "Create a beautiful web site for todogo"
    3 [2019-Sep-07] ○ : Create a beautiful web site for todogo

Each task is assigned a usage index (1,2,3) to refer to them with the
command line (no mouse to click)

============================
Listing the tasks - ``list``
============================

Then you can have a look on the whole todo list:

.. code:: shell

   $ todo list
   
    1 [2019-Sep-07] ○ : Write the documentation of todogo
    2 [2019-Sep-07] ○ : Write the unit tests of todogo
    3 [2019-Sep-07] ○ : Create a beautiful web site for todogo

   Legend: ○ todo  ▶ doing  ● done
   

The listing display for each task:

* The usage index (1,2,3, ...)
* The date of creation
* The status [todo, doing, done] (see the legend)
* The description text

=========================
Task on board - ``board``
=========================

*"This week you plan to work on the documentation (task 1) and the web
site (task 3)"*, then you can star these tasks by putting them on the
board:

.. code:: shell
   
   $ todo board -a 1,3
   Task of index 1 has been added on board
   Task of index 3 has been added on board

And list the tasks on board to focus on the actuality:

.. code:: shell

   $ todo board
   
    1 [2019-Sep-07] ○ : Write the documentation of todogo
    3 [2019-Sep-07] ○ : Create a beautiful web site for todogo
    
   Legend: ○ todo  ▶ doing  ● done
   
============================
Task life cycle - ``status``
============================

*"You start by writing some documentation and want to point that the
task is in progress"*, then you specify that you jump to the next
status of this task 1 (the status *doing*):

.. code:: shell

   $ todo status -n 1
    1 [2019-Sep-07] ▶ : Write the documentation of todogo

Then the board indicates:

.. code:: shell

   $ todo board

    1 [2019-Sep-07] ▶ : Write the documentation of todogo
    3 [2019-Sep-07] ○ : Create a beautiful web site for todogo
    
   Legend: ○ todo  ▶ doing  ● done
   
Note that todogo defines three possible status:

``○ todo``: the task is registered and is waiting to be done

``▶ doing``: the task is started and is in progress

``● done``: the task is achieved
