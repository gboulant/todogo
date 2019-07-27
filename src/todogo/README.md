This folder contains the todogo application source files.

The todogo application is composed of executable programs (prog
folder) depending on several packages/libraries (other folders).

The prog folder may contain several executable programs, one per
sub-folder whose directory name is the name of the executable program
created by the command ``go build/install`` (see Makefile).

The packages ``mypkg`` should be imported using the path::

    import "todogo/mypkg"
