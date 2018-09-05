# flowcat
Get updated task lists during development based on comments in your code.

PREREQUISITES:

python
argopt python module

SETUP:

If you do not have argopt installed you can install it with <code>pip install argopt</code>

ABOUT:

Flowcat is a python program that parses the working directory or files of your development project and returns a list of tasks that need to be completed. This works by using a regex to parse comments in your code.

As an example if you leave a comment in your code #@todo Create function to sanitize variables
flowcat lets you create a list of all the comments you have left in your files about what needs to be done and will return a list of each comment starting with #@todo.Flowcat will parse recursively through a directory you specify or let you specify a single file.

While #@todo is the default, it can be replaced by any regular expression by specifying it as the -m argument.

Works great in team development environments as well! No need to track your outstanding tasks in different software or seperate files. With flowcat note the task that needs to be done while you are writing the code. 

If multiple people are working on a file simply note the name of the person in the comment to assign the task to them and let them use flowcat to match the tasks that belong to them.

As an example:
<code>
#@marvin Fix the sanitation of variables
SOME CODE
#@Arthur Create a new menu item for the babel fish on our site.
#@marvin Create happy() function to define things that make you happy.
</code>
Now when we run taskcat we can specify the regex in the argument as <code>-m "#@marvin"</code> to get all of Mavin's tasks and <code>-m "#@Arthur"</code> to get a list of Arthur's tasks.

OPTIONS:

Usage:
  todolist [-h] (-f folder) [-o outfile] [-l] [-m match]

Options:
  -h --help     Show this screen.
  -l  Show line numbers
  -f /path/to/folder or file
  -o /path/to/outputfile  File will be overwritten if it exists
  -m (OPTIONAL) regex to match a todo item for ex -m "#@note" matches anything after #@note. Defaults to "#@todo\ "
  
  
  
  
