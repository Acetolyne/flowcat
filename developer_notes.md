This project uses a custom lexer developed to be used with flowcat available at 

github.com/Acetolyne/commentlex

Before testing new functionality for flowcat make sure you have the latest lexer version by running the below

	export GOPROXY=direct
	go get github.com/Acetolyne/commentlex

Running flowcat init will create a new .flowcat settings file in the users home directory
This settings file will include regex patterns for excluding files from the flowcat analysis
Flowcat will by default exclude the output file when you use the -o flag


Readme is updated on push to any branch

Binaries are only autobuilt after a PR is committed to Master branch


