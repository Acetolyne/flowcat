# flowcat 2.0.0

## Get updated task lists during development based on comments in your code

### PREREQUISITES

None

### SETUP

To install flowcat run ```git clone https://github.com/Acetolyne/flowcat.git```
put the flowcat binary for your OS in one of your userpaths such as /usr/local/bin/, /bin/, /sbin/ or another path

Flowcat is currently only tested on Linux systems, now that flowcat is written in GoLang I am looking at porting it to Windows systems as well. While not tested yet it should currently work with MacOS as well but feedback is required.

## VSCode autorun option

A useful workflow is to have flowcat regenerate your todo list when you save a project. Using the RunOnSave extension by pucelle this can be accomplished, you have to setup a separate setting for each project folder but it can really help you track your tasks during your projects lifetime. Below is an example of how I have mine setup. these setting need to go in the settings.json file in VSCode found by navigating to File > Preferences > Settings then clicking on Run On Sve under extensions, which is visible after installing the extension, then click edit in settings.json and modify the file to have something similar to the below example replacing the paths to your project folders. Add as many project folders as you need, this also allows for different regex matching for different programming languages or even creation of multiple separated task lists.

```
"runOnSave.commands": [
    
        {
            "match": "/home/acetolyne/Project1/*",
            "command": "flowcat -f /home/acetolyne/Project1/ -o /home/acetolyne/Project1/todo -l -m '//@todo'",
            "runIn": "terminal",
            "runningStatusMessage": "Updating task list",
            "finishStatusMessage": "Task list updated"
        },
        {
            "match": "/home/acetolyne/Project2/*",
            "command": "flowcat -f /home/acetolyne/Project2/ -o /home/acetolyne/Project2/todo -l -m '//@todo'",
            "runIn": "terminal",
            "runningStatusMessage": "Updating task list",
            "finishStatusMessage": "Task list updated"
        },
    ],
```

### ABOUT

Flowcat is a GoLang program that parses the working directory or files of your development project and returns a list of tasks that need to be completed. This works by using a regex to parse comments in your code.

As an example if you leave a comment in your code //@todo Create function to sanitize variables
flowcat lets you create a list of all the comments you have left in your files about what needs to be done and will return a list of each comment starting with #@todo.Flowcat will parse recursively through a directory you specify or let you specify a single file.

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
