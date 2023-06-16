# easygpt

An app that uses the chatgpt API to process text files in bulk.  
It can be used for translation, adding comments to source code, and more, depending on how you use it.  

## Install

```sh
go install github.com/xcd0/easygpt@latest
```

## Simple Example (Getting output by giving a string to AI)

Set your OpenAI API key to the environment variable `OPENAI_API_KEY`.  
The API key below is an invalid API key for the sake of example, so issue your own valid API key and set it.  
You can issue an API key from https://platform.openai.com/account/api-keys.  
Set the issued API key to `.bashrc` or similar.

```sh
$ export OPENAI_API_KEY=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
$ echo export OPENAI_API_KEY=$OPENAI_API_KEY >> .bashrc
```

Then install easygpt and run it.

```sh
$ go install github.com/xcd0/easygpt@latest
$ easygpt -i "Say this is test."
this is test.
```

## Batch Processing of Multiple Files

Process multiple text files in bulk with AI.

### Preparation

Generate a template for the setting file and edit it to create the setting file.

1. Generate a template.

```sh
$ easygpt --create-setting
```

When you run it like this, `eagygpt.hjson` will be generated in the current directory.  
You can also use `easygpt -c`.

2. Edit the setting file.
* Set the API key
	* It is recommended to set it to the environment variable `OPENAI_API_KEY`. In this case, you do not need to write it in the setting file.
	* You can also write the API key in the `apikey` section of the setting file.  
		Example: `apikey: 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`  
		The above API key is an invalid API key, so issue your own.  
		You can issue an API key from https://platform.openai.com/account/api-keys.

* Write the string you want to give before the input text file in the `prompt` section of the setting file.  
	For example, if you want to translate an English text file, write `prompt: Translate the following.`  
	If you want to write multiple lines, use `\n` to insert a line break or write it in the following heredoc format.  
	```
	prompt:  
	'''  
	You can write multiple lines like this.  
	'''  
	```

For other arguments, read the comments in the setting file and write them as needed.

3. Place the setting file.

The setting file is searched according to the following specifications.  
Basically, you should place easygpt.hjson in the current directory.

The default names for the setting file are:
* easygpt.hjson
* .easygpt.hjson
* .easygpt

Other than these names, you can set the setting file with the `--setting` argument.

The default locations for the setting file are:
* Current directory
* Home directory
* The same directory as the executable file

The search for the setting file is performed in this order.  
Therefore, if it exists in the current directory, the setting file in the home directory will be ignored.

## Usage 1: Drag and Drop

Usage on the GUI. Select files or folders and drag and drop them onto the executable file.  
All the given files and files included in the given folders are processed by AI in bulk.  
The output is saved with the name of the input file plus `_easygpt_output`.

If a setting file exists, you can use it as follows:
1. Drag and drop the text file or directory you want to throw to gpt onto the easygpt executable file.
1. The processing result is output with the name of the input file plus `_easygpt_output`.

When using this method, the following settings in the setting file will be ignored:
* input-dir
* output-dir
* extension

## Usage 2: Execution from the Command Line

As an example of a simple way to run, if the API key is set in the environment variable or the setting file, you can use it as follows:

```sh
$ ./easygpt -i Introduce yourself.
Hello, I am an AI. I am a natural language processing model developed by OpenAI. My purpose is to provide the best answers and responses when users ask questions or make requests. I have information on various topics and can also make grammar and style corrections. How can I assist you?
```

## Usage 3: Execution from the Command Line - Processing a Group of Files

If command line arguments are set, the values in the setting file are overridden by the values in the command line arguments.  
For the values set in the default path setting file, you can omit the arguments unless you want to change them.  
It is recommended to use the setting file because there are many arguments that can be set.

There are flag arguments in the form of `--key value` and non-flag arguments in the command line arguments.  
* Flag arguments have priorities, and if a specific argument is used, other arguments may be ignored.  
  For example, if `--create-setting` exists, other arguments are ignored and the program generates the template for the setting file and exits.  
  Unless you want to leave the value empty, you do not need to enclose it in "".  
* Non-flag arguments are treated as input files.  
  If a directory is specified, all files in the directory are treated as input files recursively.

The arguments can be in any order.  
For details on the arguments, refer to the help text output by `./easygpt --help`.

## Open Source Software Libraries Used

The license for each open source software library is stored in `./licenses`.

* Argument parsing and help generation
	* [alexflint/go-arg](https://github.com/alexflint/go-arg)
		* BSD-2-Clause license
* Generation of random strings (UUID is not necessary)
	* [google/uuid](https://github.com/google/uuid)
		* BSD-3-Clause license
* Stack trace
	* [pkg/errors](https://github.com/pkg/errors)
		* BSD 2-Clause "Simplified" License
* Setting file
	* [hjson/hjson-go](https://github.com/hjson/hjson-go)
		* MIT License
* Getting the Longest Common Subsequence (LCS)
	* [go.pkgs](github.com/cloudengio/go.pkgs)
		* Apache License 2.0
* Golang itself
	* [golang/go](https://github.com/golang/go)
		* BSD 3-Clause "New" or "Revised" License
* Getting the license files
	* [google/go-licenses](https://github.com/google/go-licenses)
		* Apache License 2.0
* Executable file compression
	* [upx/upx](https://github.com/upx/upx)
		* GPL2+ or UPX LICENSE
			* https://upx.github.io/upx-license.html
		* The program is distributed without modification by compressing it with upx during the build process. The source code of upx is not included in the program.

## LICENSE

MIT License