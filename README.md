# go-lava-bomb

This volcano is one ticking bomb!

--------------------------------------------------------

# Requirements

Create a CLI Lava Bomb application that will rumble/print the following noises/messages at the following intervals to stdout:

- "rumble" every second
- "RUMBLE" every minute
- "LAVAOVERFLOW" every hour

Only one value should be printed in a given second, e.g. on the hour, only "LAVAOVERFLOW" will be printed.

The Lava Bomb application should run for four hours and then exit.

A mechanism should exist for the user to alter any of the printed values while the program is running, i.e. after the application has run for 15 minutes, it should be possible, without stopping the program, to change what is printed either per second or per minute, per hour or possibly all noises can change. 

# Implementation

To be able to dynamically change displayed messages, application must monitor some source from which it can receive new values during the runtime. There are various ways to achieve this, for example watching a local file for content update or using a pub/sub messaging system. Here the simplest, file based, approach was used.

Entry point to the application is main.go, where a new volcano object is created, which then starts erupting.
The application can be started with a -file flag to specify a custom file, which contains custom noises/messages, and which, on updating and saving, will cause a change in the printouts to the stdout. If no flag is provided, the default file is config.json.

The Erupt method starts three goroutines:
    - rumble() is responsible for time keeping, ieoutput at prescribed intervals
    - monitor() is responsible for checking the source of custom changes, ie a file in this case
    - print() prints out a message when it receives one from a channel (to which it is sent by the rumble goroutine)

sending on done channel and closing the noiseCh channel, to which messages are sent, is used to gracefully finish and close the application.


The functionality itself is contained in the volcano package.

# Usage

- go run main.go

- go run main.go --file=json-file-path

        default file, if the config flag is not provided, is config.json

- make dev

  CUSTOM_FILE=/data/custom.json  make run

        custom config file needs to start with "/data/" followed by its relative position to the current working directory

# Testing

Test of the rumble method exceeds the default 30s limit timeout, so it is necessary to provide a custom timeout:

    go test -timeout 180s github.com/tamarakaufler/go-lava-bomb/volcano
    
    go test -coverprofile=/tmp/go-code-cover -timeout 180s github.com/tamarakaufler/go-lava-bomb/volcano

    go test -timeout 180s github.com/tamarakaufler/go-lava-bomb/volcano -run Test_volcano_rumble

Unit tests provide 60.0% coverage.
