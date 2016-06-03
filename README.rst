=======
seekret
=======

Go library and command line to seek for secrets on various sources.


************
Command Line
************

Synopsis
========

::

    NAME:
       seekret - seek for secrets on various sources.   

    USAGE:
       seekret [global options] command [command options] [arguments...]
       
    VERSION:
       0.0.1
       
    AUTHOR(S):
       Albert Puigsech Galicia <albert@puigsech.com> 
       
    COMMANDS:
       seek:
         git    seek for seecrets on a git repository.
         dir    seek for seecrets on a directory.   

    GLOBAL OPTIONS:
       --exception FILE, -x FILE    load exceptions from FILE.
       --rules PATH         PATH with rules. [$SEEKRET_RULES_PATH] 
       --format value, -f value specify the output format. (default: "human")
       --help, -h           show help
       --version, -v        print the version



Description
===========

``seekret`` inspect different sources (files into a directory or git repositories)
to seek for secrets. It can be used to prevent that secrets are published in exposed
locations.


Installing seekret
==================

``seekret`` can be directly installed by using go get.

::

    go get github.com/apuigsech/seekret/cmd/seekret


Options
========

General Options
~~~~~~~~~~~~~~~

``-x, --exception``

``--rules``

``-f, --format``


Options for Git
~~~~~~~~~~~~~~~

``-c, --count``


Options for Dir
~~~~~~~~~~~~~~~

``-r, --recursive``

``-h, --hidden``


Examples
========

Scan all files from all commits in a local repo::

    seekret git /path/to/repo

Scan all files from all commits in a remote repo::

    seekret git http://github.com/apuigsech/seekret-exposed

Scan all files from the last commit in a local repo::

    seekret git --count 1 /path/to/repo

Scan all files (including hidden) in a local folder::

    seekret dir --recursive --hidden /path/to/dir



*******
Library
*******

TBD
