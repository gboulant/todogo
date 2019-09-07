:title: TODOGO - Userguide
:author: Guillaume Boulant (EDF/R&D)
:date: Sept. 2019
:description: Quick start guide for TODOGO

-------------

.. raw:: html

   <div align="center" style="padding-top: 20%; padding-left:20%; padding-right:20%">
   <h1 style="margin-left: 0; margin-right: 0">01 - Userguide</h1>
   <p style="margin-left: 0; margin-right: 0">
   </p>
   </div>

==============
Creating tasks
==============

The todo application can help you to create and manage a personal todo
list. For creating a new task::

   $ todo new -t "Write the documentation of todogo"
    1 [2019-07-27]  todo: Write the documentation of todogo

The todo task is created with usage index 1. You may complete your
todo list:

.. code:: shell

   $ todo new -t "Write the unit tests of todogo"
    2 [2019-07-27]  todo: Write the unit tests of todogo

   $ todo new -t "Create a beautiful web site for todogo"
    3 [2019-07-27]  todo: Create a beautiful web site for todogo

Each task is assigned a usage index (1,2,3) to refer to them in the
folowing use cases.

=================
Listing the tasks
=================

Then you can have a look on the whole todo list:

.. code:: shell

  guillaume@mercure:~$ todo list
  
   1 [2019-Aug-15] ▶ : Write the documentation of todogo
   2 [2019-Aug-15] ▶ : Write the unit tests of todogo
   3 [2019-Aug-15] ○ : Create a beautiful web site for todogo
   4 [2019-Aug-15] ● : Develop a tree representation of parent-children tasks
   5 [2019-Aug-15] ○ : Push a clone of the repository on github
   6 [2019-Aug-15] ● : Create a dockerfile of the todogo application

  Legend: ○ todo  ▶ doing  ● done

-------------

.. raw:: html

   <div align="center" style="padding-top: 20%; padding-left:20%; padding-right:20%">
   <h1 style="margin-left: 0; margin-right: 0">02 - Getting started</h1>
   <p style="margin-left: 0; margin-right: 0">
   </p>
   </div>

-------------

.. raw:: html

   <div align="center" style="padding-top: 20%; padding-left:20%; padding-right:20%">
   <h1 style="margin-left: 0; margin-right: 0">03 - Technical design</h1>
   <p style="margin-left: 0; margin-right: 0">
   </p>
   </div>

