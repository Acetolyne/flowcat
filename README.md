# flowcat 3.0.0


## Flowcat helps developers bring their todo tasks forward to the user so you don't need to waste time looking for them. It provides a clear overview of your tasks and allows you to stay focused on your work. You can also output to a file for planning purposes or to separate tasks into different categories

<a href="https://goreportcard.com/report/github.com/Acetolyne/flowcat" target="_blank"><img src="https://goreportcard.com/badge/github.com/Acetolyne/flowcat?style=flat&logo=none" alt="go report" /></a>

### PREREQUISITES

None

### ABOUT

Flowcat is a GoLang program that parses the working directory or files of your development project and returns a list of tasks that need to be completed. This works by looking thru all the comments in your files recursively and returning any comments that match a specified string.

As an example if you leave a comment in your code 
``//TODO Create function to sanitize variables``
flowcat lets you create a list of all the comments you have left in your files about what needs to be done and will return a list of each comment starting with @todo.Flowcat will parse recursively through a directory you specify.

While ``TODO`` is the default, it can be replaced by any regular expression by specifying it as the -m argument. You can also set a new default for your user by first running flowcat init then editing the file at ~/.flowcat/config

Works great in team development environments as well! No need to track your outstanding tasks in different software or seperate files. With flowcat you can note the task that needs to be done while you are writing the code meaning you don't have to shift focus and come back to your development environment.

If multiple people are working on a file simply note the name of the person in the comment to assign the task to them and let them use flowcat to match the tasks that belong to them.

As an example:

```golang
//@marvin Fix the sanitation of variables
SOME CODE
//@Arthur Create a new menu item for the babel fish on our site.
//@marvin Create happy() function to define things that make you happy.
```

Now when we run flowcat we can specify the regex in the argument as ``-m "@marvin"`` to get all of Mavin's tasks and ``-m "@Arthur"`` to get a list of Arthur's tasks.

If you do not run flowcat init after you have downloaded it then there may be some issues as flowcat will not ignore any files or folders and will recurse into .git and other hidden directories.


### OPTIONS
```text
Usage:
  flowcat [-h] (-f folder) [-o outfile] [-l] [-m match]

Options for Flowcat:
-f string
    The project top level directory, where flowcat should start recursing from. (default '.' Current Directory)
-l    If line numbers should be shown with todo items in output.
-m string
    The string to match to do items on. (default '//@todo')
-o string
    Optional output file to dump results to, note output will still be shown on terminal.
-h
    This help menu
```

### SETUP

#### Installation
Download the appropriate version for your system below then put the flowcat binary for your OS in one of your userpaths such as /usr/local/bin/, /bin/, /sbin/ or another path

[![functionality](https://github.com/Acetolyne/flowcat/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/Acetolyne/flowcat/actions/workflows/test.yml)

[![build](https://github.com/Acetolyne/flowcat/actions/workflows/build.yml/badge.svg?branch=master)](https://github.com/Acetolyne/flowcat/actions/workflows/build.yml)


Flowcat is currently only tested on Linux systems, now that flowcat is written in GoLang I am looking at porting it to Windows systems as well. While not tested yet it should currently work with MacOS as well but feedback is required.

#### Build From Source

To install from source clone the repository with the command ```git clone https://github.com/Acetolyne/flowcat```
build the binary for your system with ```go build *.go -o flowcat```
Now move the binary to a system path such as /usr/bin/

#### download binary
alternatively if you don't want to setup a go environment to build from source you may use one of the pre-compiled binaries that are aailable as releases on the github page https://github.com/Acetolyne/flowcat.

#### Settings initialization
running ```flowcat init``` will allow you to make a settings file for your user so that you can set the default -m argument as well as a list of regex for files to ignore. Once running flowcat init flowcat will confirm the file was written to ~/.flowcat/config, it is suggested to edit the settings file to include the files or regex you would like to ignore when flowcat runs.

#### Regex
The regex used for matching against ignore files and to match lines via the -m argument is described in the documentation at https://github.com/google/re2/wiki/Syntax and uses the MatchString method.

Certain characters will need to be escaped with a backslash ```\``` including the backslash character itself ```\``` becomes ```\\```

The syntax of the regular expressions accepted is the same general syntax used by Perl, Python, and other languages. More precisely, it is the syntax accepted by RE2 and described at https://golang.org/s/re2syntax, except for \C. For an overview of the syntax, run ```go doc regexp/syntax```




#### VSCode autorun option

A useful workflow is to have flowcat regenerate your todo list when you save a file. Using the RunOnSave extension by pucelle this can be accomplished, you have to setup a separate setting for each project folder but it can really help you track your tasks during your projects lifetime. Below is an example of how I have mine setup. these setting need to go in the settings.json file in VSCode found by navigating to File > Preferences > Settings then clicking on Run On Save under extensions, which is visible after installing the extension, then click edit in settings.json and modify the file to have something similar to the below example replacing the paths to your project folders. Add as many project folders as you need, this also allows for different regex matching for different programming languages or even creation of multiple separated task lists.

```json
"runOnSave.commands": [
    
        {
            "match": "/home/acetolyne/Project1/*",
            "command": "flowcat -f /home/acetolyne/Project1/ -o /home/acetolyne/Project1/todo -l -m '@todo'",
            "runIn": "terminal",
            "runningStatusMessage": "Updating task list",
            "finishStatusMessage": "Task list updated"
        },
        {
            "match": "/home/acetolyne/Project2/*",
            "command": "flowcat -f /home/acetolyne/Project2/ -o /home/acetolyne/Project2/todo -l -m '@todo'",
            "runIn": "terminal",
            "runningStatusMessage": "Updating task list",
            "finishStatusMessage": "Task list updated"
        },
    ],
```
### Supported filetypes
flowcat currently supports the following filetypes additionally files with no extensions use the basic // comment style and /*  */ comment style for multiline comments.
To have additional filetypes added to flowcat please open an issue on GitHub for the lexer that is used by flowcat at github.com/Acetolyne/commentlex

#### Supported Filetypes <!--Everything below this line is autogenerated do not edit -->

```text
.c
.class
.cpp
.go
.gohtml
.h
.html
.jar
.java
.js
.jsp
.lua
.md
.php
.py
.rb
.rs
.sh
.tmpl
```