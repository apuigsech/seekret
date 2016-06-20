=======
seekret
=======

Go library and command line to seek for secrets on various sources.


************
Command Line
************

Description
===========

``seekret`` inspect different sources (files into a directory or git 
repositories) to seek for secrets. It can be used to prevent that secrets are
published in exposed locations.


Installing seekret
==================

``seekret`` can be directly installed by using go get.

::

    go get github.com/apuigsech/seekret/cmd/seekret


Usage
=====

General Options
~~~~~~~~~~~~~~~

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
       --known FILE, -k FILE  load known secrets from FILE.
       --workers value, -w value  number of workers used for the inspection (default: 4)
       --help, -h           show help
       --version, -v        print the version


``-x, --exception``

``--rules``

``-f, --format``

``-k, --known``

``-w, --workers``


Options for Git
~~~~~~~~~~~~~~~

::

    NAME:
       seekret git - seek for seecrets on a git repository.

    USAGE:
       seekret git [command options] [arguments...]

    CATEGORY:
       seek

    OPTIONS:
       --count value, -c value  (default: 0)
   

``-c, --count``


Options for Dir
~~~~~~~~~~~~~~~

::

    NAME:
       seekret dir - seek for seecrets on a directory.  

    USAGE:
       seekret dir [command options] [arguments...] 

    CATEGORY:
       seek 

    OPTIONS:
       --recursive, -r  
       --hidden


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


Hands-On
========

The repository seekret-secrets is prepare to test seekret, and can be used to
perform the following hands-on examples:

1. Inspect remote git repository::

    seekret --rules $GOPATH/src/github.com/apuigsech/seekret/rules/ git https://github.com/apuigsech/seekret-secrets.git 

2. Inspect local got repository::

    git clone https://github.com/apuigsech/seekret-secrets.git /tmp/seekret-secrets
    seekret --rules $GOPATH/src/github.com/apuigsech/seekret/rules/ git /tmp/seekret-secrets

3. Inspect only the last 2 commits::

    seekret --rules $GOPATH/src/github.com/apuigsech/seekret/rules/ git -c 2 /tmp/seekret-secrets

4. Inspect with exceptions::

    seekret --rules $GOPATH/src/github.com/apuigsech/seekret/rules/ -x /tmp/seekret-secrets/.exception_1 git /tmp/seekret-secrets


*******
Library
*******

Importing seekret Library
=========================

::

    import seekret "github.com/apuigsech/seekret/lib"


Init Seekret context
====================

::

    s := seekret.NewSeekret()


Loading Rules
=============

::

    s.LoadRulesFromPath("/path/to/main/rues:/path/to/other/rules:/path/to/more/rules")

::

    s.LoadRulesFromDir("/path/to/rules")


::

    s.LoadRulesFromFile("/path/to/file.rule")


Loading Objects
===============

::

    opts := map[string]interface{} {
      "hidden": true,
      "recursive": false,
    }
    s.LoadObjects("dir", "/path/to/inspect", opts)


::

    opts := map[string]interface{} {
      "count": 10,
    }
    s.LoadObjects("dir", "/repo/to/inspect", opts)


Loading Exceptions
==================

::

    s.LoadExceptionsFromFile("/path/to/exceptions/file")



Inspect
=======

::

    s.Inspect(5)



Get Inspect Results
===================

::

    secretsList := s.ListSecrets()



*****
Rules
*****

Secret identification is performed by using a set of rules specified on the
rules files. Those files, with '.rule' extension are defined by using YAML 
following this format: 

::

    rulename:
      match: [regexp]
      unmatch:
        - [regexp]
        - [regexp]
        - ...

For the contents of a file is considered a secret, it must comply with the
'match' regexp and not comply ANY of the 'unmatch' reg rule and comply match 
ANY of the unmatch.


**********
Exceptions
**********

Exceptions determine conditions under which content should not be considered
a secret. The exceptions are specified by using a YAML file that follows this
format:

::
    
      ...
    -
      rule: [rulename]
      object: [regexp]
      line: [linenumber]
      content: [regexp]
    -
      ...


The conditions are optional, so it is not necessary to specify them all, but
for a content deemed exception must meet all the specified conditions.

The meaning of the various conditions explained:

``rule``
Contains the name of the rule.

``object``
Contains a regexp that should match the object name (usually the filename).

``line``
Contains the line number into the object.

``content``
Contains a regexp that should match the content.