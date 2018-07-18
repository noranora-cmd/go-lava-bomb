# Requirements

Create a CLI Lava Bomb application that will rumble/print the following noises/messages at the following intervals to stdout:

- "rumble" every second
- "RUMBLE" every minute
- "LAVAOVERFLOW" every hour

Only one value should be printed in a given second, e.g. on the hour, only "LAVAOVERFLOW" should be printed.

The Lava Bomb application should run for four hours and then exit.

A mechanism should exist for the user to alter any of the printed values while the program is running, i.e. after the clock has run for 10 minutes I should, without stopping the program, be able to change it so that it stops printing "tick" every second and starts printing "quack" instead. 

# Implementation

To be able to dynamically change displayed messages, application must monitor some source from which it can receive new values during the runtime. There are various ways to achieve this, for example watching a local file for content update or using a pub/sub messaging system. Here the simplest, file based, approach was used.

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