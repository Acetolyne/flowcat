# flowcat
<b>Get updated task lists during development based on comments in your code.</b>

PREREQUISITES:

python >= 2.7<br>
argopt python module

SETUP:

To install flowcat run <code>git clone https://github.com/Acetolyne/flowcat.git</code><br>
If you do not have argopt installed you can install it with <code>pip install argopt</code><br>
You can now run flowcat from the directory you downloaded it to by running <code>./flowcat -f FOLDER/FILE OTEHR OPTIONS</code><br>
To ba able to run it from anywhere create a symbolic link in your /usr/bin or /usr/share/bin directory that points to the location of the flowcat executable.


ABOUT:

Flowcat is a python program that parses the working directory or files of your development project and returns a list of tasks that need to be completed. This works by using a regex to parse comments in your code.

As an example if you leave a comment in your code #@todo Create function to sanitize variables
flowcat lets you create a list of all the comments you have left in your files about what needs to be done and will return a list of each comment starting with #@todo.Flowcat will parse recursively through a directory you specify or let you specify a single file.

While #@todo is the default, it can be replaced by any regular expression by specifying it as the -m argument.

Works great in team development environments as well! No need to track your outstanding tasks in different software or seperate files. With flowcat note the task that needs to be done while you are writing the code. 

If multiple people are working on a file simply note the name of the person in the comment to assign the task to them and let them use flowcat to match the tasks that belong to them.

As an example:
<pre>
<code>
#@marvin Fix the sanitation of variables
SOME CODE
#@Arthur Create a new menu item for the babel fish on our site.
#@marvin Create happy() function to define things that make you happy.
</code></pre>
Now when we run taskcat we can specify the regex in the argument as <code>-m "#@marvin"</code> to get all of Mavin's tasks and <code>-m "#@Arthur"</code> to get a list of Arthur's tasks. All tasks should be a comment in your code in PHP your line should start with <code>\\</code> in python and other laguages it will need to start with a <code>#</code>

OPTIONS:

Usage:
  todolist [-h] (-f folder) [-o outfile] [-l] [-m match]
<pre>
Options:
  -h --help     Show this screen.
  -l  Show line numbers
  -f /path/to/folder or file
  -o /path/to/outputfile  File will be overwritten if it exists
  -m regex to match a todo item for example -m "#@note" matches anything after #@note. Defaults to "#@todo\ "
</pre>
