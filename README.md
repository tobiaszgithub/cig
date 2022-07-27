# CLI for Cloud Integration.
Using this tool it is possible to manage and query integration artifacts
of design time and runtime.

Usage:<br>
&ensp;cig [command]

Available Commands:
- completion -      Generate the autocompletion script for the specified shell
- flow -            Command related to the processing of an integration flow
- generate-config - Generate config file
- help -            Help about any command
- package -         Command related to the processing of integration packages

Flags:<br>
&ensp;-h, --help&ensp;&ensp;help for cig

Use "cig [command] --help" for more information about a command.


## cig flow
Command related to the processing of an integration flow.

Usage:<br>
&ensp;cig flow [command]

Available Commands:
- copy -             Copy an integration flow
- create -           Create or upload an integration flow
- deploy -           Deploy an integration flow
- describe-configs - Get configurations of an integration flow by Id and version
- download -         Download an integration flow as zip file
- inspect -          Get integration flow by id and version
- transport -        Transport an integration flow between systems
- update -           Update an integration flow
- update-configs -   Update configuration parameters of an integration flow

Flags:<br>
&ensp;-h, --help&ensp;&ensp;help for flow

## cig generate-config
Generate configuration file. This file is nessesary for the operation
of the cig tool. Configuration file should be placed in working directory or userhome/.cig/ directory

Usage:<br>
&ensp;cig generate-config [flags]

Flags:<br>
&ensp;-h, --help&ensp;&ensp;help for generate-config<br>
&ensp;-o, --output-file&ensp;&ensp;string&ensp;&ensp;The output file with empty configuration parameters that will be created (default "config.json")

## cig package
Command related to the processing of integration packages

Usage:<br>
&ensp;cig package [flags]<br>
&ensp;cig package [command]

Aliases:<br>
&ensp;package, ls, p

Available Commands:
- download -    Download integration package by ID
- inspect -     Get integration package by ID
- ls -          Get all integration packages as list or get all integration flow of the package

Flags:<br>
&ensp;-h, --help&ensp;&ensp;help for package