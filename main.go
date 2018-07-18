package main

import (
	"flag"
	"log"
	"time"

	"github.com/tamarakaufler/go-lava-bomb/volcano"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "file", "config.json", "Path of a json file through which volcano noises can be dynamically changed")
}

func main() {
	flag.Parse()

	c, err := volcano.New(time.Duration(240 * time.Minute))
	if err != nil {
		log.Fatal(err)
	}
	c.Erupt(configFile)
}
