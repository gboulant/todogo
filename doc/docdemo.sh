#!/bin/sh

todo config -n "demo"

todo new -t "Write the documentation of todogo"
todo new -t "Write the unit tests of todogo"
todo new -t "Create a beautiful web site for todogo"

todo board -a 1,3

todo status -n 1
todo board

todo status -n 1
todo board -r 1

todo archive -a 1
todo list
todo archive


todo new -t "Develop a tree representation of parent-children tasks"
todo new -t "Push a clone of the repository on github"
todo new -t "Create a dockerfile of the todogo application"

todo status -n 1,2,4,6
todo status -n 4,6

todo list


todo new -t "Write the conceptual design of the dingo application"
todo new -t "Setup the technical environment for the dingo application"
todo new -t "Phone IT center to get a new PC"
todo new -t "Book an hotel for the workshop in Melun"
todo new -t "Write a prototype of dingo to validate the design"
todo new -t "Write a project proposition to get a budget for dingo"



