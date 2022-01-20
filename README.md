# flowcat 2.0.0

## Get updated task lists during development based on comments in your code

<a href="https://goreportcard.com/report/github.com/Acetolyne/flowcat" target="_blank"><img src="https://goreportcard.com/badge/github.com/Acetolyne/flowcat?style=flat&logo=none" alt="go report" /></a>

### PREREQUISITES

None

### ABOUT

Flowcat is a GoLang program that parses the working directory or files of your development project and returns a list of tasks that need to be completed. This works by using a regex to parse comments in your code.

As an example if you leave a comment in your code //@todo Create function to sanitize variables
flowcat lets you create a list of all the comments you have left in your files about what needs to be done and will return a list of each comment starting with //@todo.Flowcat will parse recursively through a directory you specify.

While //@todo is the default, it can be replaced by any regular expression by specifying it as the -m argument.

Works great in team development environments as well! No need to track your outstanding tasks in different software or seperate files. With flowcat note the task that needs to be done while you are writing the code.

If multiple people are working on a file simply note the name of the person in the comment to assign the task to them and let them use flowcat to match the tasks that belong to them.

As an example:

```golang
//@marvin Fix the sanitation of variables
SOME CODE
//@Arthur Create a new menu item for the babel fish on our site.
//@marvin Create happy() function to define things that make you happy.
```

Now when we run flowcat we can specify the regex in the argument as ```-m "//@marvin"``` to get all of Mavin's tasks and ```-m "//@Arthur"``` to get a list of Arthur's tasks. All tasks should be a comment in your code in PHP your line should start with ```\\``` in python and other laguages it will need to start with a ```#```

Note that flowcat will not recurse into hidden directories anymore, this caused a bug when parsing folders that had a .git directory.
If you need to you can specify those directories in the -f argument to get the todo items from there but there should be no need to parse hidden directories.
If there is a need please put in a feature request on GitHub.

### OPTIONS
```
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

[![linux-amd64 build and functionality](https://github.com/Acetolyne/flowcat/actions/workflows/linux-amd64%20build%20and%20functionality.yml/badge.svg?branch=master)](https://github.com/Acetolyne/flowcat/actions/workflows/linux-amd64%20build%20and%20functionality.yml) https://www.github.com/Acetolyne/flowcat/bin/flowcat-linux-amd64/flowcat

[![linux-386 build and functionality](https://github.com/Acetolyne/flowcat/actions/workflows/linux-386%20build%20and%20functionality.yml/badge.svg?branch=master)](https://github.com/Acetolyne/flowcat/actions/workflows/linux-386%20build%20and%20functionality.yml) https://www.github.com/Acetolyne/flowcat/bin/flowcat-linux-386/flowcat



Flowcat is currently only tested on Linux systems, now that flowcat is written in GoLang I am looking at porting it to Windows systems as well. While not tested yet it should currently work with MacOS as well but feedback is required.

#### Installation from source

To install from source clone the repository with the command ```git clone https://github.com/Acetolyne/flowcat```
build the binary for your system with ```go build -o flowcat```
Now move the binary to a system path such as /usr/bin/

#### Project initialization
running ```flowcat init``` will allow you to make a settings file for your project so you don't need to pass arguments if it is being run against the currect folder. This is the prefered way as it will also allow you to specify regex patterns for files to ignore. Once running flowcat init flowcat will ask you for the settings you would like to use and these will be written to a settings file called .flowcat, it is suggested to edit the .flowcat file to include the files or regex you would like to ignore will flowcat runs.

#### Regex
The regex used for matching against ignore files and to match lines via the -m argument is described in the documentation at https://github.com/google/re2/wiki/Syntax and uses the MatchString method.

Certain characters will need to be escaped with a backslash ```\``` including the backslash character itself ```\``` becomes ```\\```

The syntax of the regular expressions accepted is the same general syntax used by Perl, Python, and other languages. More precisely, it is the syntax accepted by RE2 and described at https://golang.org/s/re2syntax, except for \C. For an overview of the syntax, run ```go doc regexp/syntax```




#### VSCode autorun option

A useful workflow is to have flowcat regenerate your todo list when you save a file. Using the RunOnSave extension by pucelle this can be accomplished, you have to setup a separate setting for each project folder but it can really help you track your tasks during your projects lifetime. Below is an example of how I have mine setup. these setting need to go in the settings.json file in VSCode found by navigating to File > Preferences > Settings then clicking on Run On Save under extensions, which is visible after installing the extension, then click edit in settings.json and modify the file to have something similar to the below example replacing the paths to your project folders. Add as many project folders as you need, this also allows for different regex matching for different programming languages or even creation of multiple separated task lists.

```
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

