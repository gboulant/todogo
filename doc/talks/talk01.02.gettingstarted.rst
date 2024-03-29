
-------------

.. raw:: html

   <div align="center" style="padding-top: 20%; padding-left:20%; padding-right:20%">
   <h1 style="margin-left: 0; margin-right: 0">02 - Getting started</h1>
   <p style="margin-left: 0; margin-right: 0"></p>
   </div>

======================================
Download the source and install todogo
======================================

The todogo application (``todo`` program) is written with the langage
go (https://golang.org). You first need to install go and basic
development tools (git, make). You are supposed here to be sudoers or
to be able to have this software programs installed on your host:

.. code:: shell

   $ sudo apt-get install git
   $ sudo apt-get install make
   $ sudo apt-get install golang

Then you can clone the source files and build the ``todo`` executable
program:

.. code:: shell

   $ git clone https://gitlab.galuma.net/guiboule/todogo.git
   $ cd todogo
   $ make
   $ make test
   $ sudo make install

This last command installs the executable program ``todo`` in the
directory ``$PREFIX/bin`` where PREFIX default to ``/usr/local``.

If you need to install todogo in another folder, replace with:

.. code:: shell
   
   $ PREFIX=/path/to/my/installdir make install

If ``/usr/local/bin`` (more generally ``$PREFIX/bin``) is in your
PATH, then you are ready to start with todogo.

===================
Docker installation
===================

This minimalist docker file can be used to create an ubuntu image
containing an installation of todogo. It is created only to show and
test the minimal software configuration required to work with todogo:

.. code:: docker

   FROM ubuntu

   RUN apt-get update && apt-get upgrade -y && \
       apt-get install -y git && \
       apt-get install -y make

   RUN apt-get install -y golang

   RUN git clone https://gitlab.galuma.net/guiboule/todogo.git && \
       cd todogo && make install


If you save the content in a file named ``Dockerfile``, then you can
build the image using the command:

.. code:: shell

   $ docker build -f Dockerfile -t todogoimg .

Then run the tests using:

.. code:: shell

   $ docker run --rm -it todogoimg make -C todogo test

And display the todogo config with:

.. code:: shell

   $ docker run --rm -it todogoimg todo config -i

Finally, you can delete the image with:

.. code:: shell

   $ docker image rm todogoimg

