#!/bin/sh

todo config -n "demo"
todo new -t "Write the documentation of todogo"
todo new -t "Write the unit tests of todogo"
todo new -t "Create a beautiful web site for todogo"
todo new -t "Develop a tree representation of parent-children tasks"
todo new -t "Push a clone of the repository on github"
todo new -t "Create a dockerfile of the todogo application"

todo status -n 1,2,4,6
todo status -n 4,6

todo list

